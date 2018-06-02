package matcher

import (
	"container/heap"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/nsqio/go-nsq"
	"github.com/quickfixgo/enum"
	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/rudeigerc/broker-gateway/service"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type MatchHandler struct {
	*Matcher
}

func (h *MatchHandler) HandleMessage(m *nsq.Message) error {
	order := model.Order{}
	order.Unmarshal(m.Body)

	switch enum.OrdType(order.OrderType) {
	case enum.OrdType_MARKET:
		service.Order{}.NewOrder(&order)
		h.NewMarketOrder(order)
	case enum.OrdType_LIMIT:
		service.Order{}.NewOrder(&order)
		h.NewLimitOrder(order)
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		service.Order{}.NewOrder(&order)
		h.NewStopOrder(order)
	case enum.OrdType_COUNTER_ORDER_SELECTION:
		h.NewCancelOrder(order)
	default:
		return fmt.Errorf("[matcher.matcher.HandleMessage] [ERROR] Invalid order type: %s", enum.OrdType(order.OrderType))
	}

	return nil
}

type Matcher struct {
	*nsq.Consumer

	etcdClient *clientv3.Client

	asksLimitOrderBook *MinHeap
	bidsLimitOrderBook *MaxHeap

	asksStopOrderBook *MaxHeap
	bidsStopOrderBook *MinHeap

	marketPrice decimal.Decimal

	executor *Executor
}

func NewMatcher() *Matcher {
	config := nsq.NewConfig()
	config.LookupdPollInterval = time.Second
	consumer, err := nsq.NewConsumer(viper.GetString("nsq.topic"), "order", config)
	if err != nil {
		log.Fatalf("[matcher.matcher.NewMatcher] [FETAL] %s", err)
	}

	m := &Matcher{
		Consumer:           consumer,
		etcdClient:         mapper.NewEtcdClient(),
		asksLimitOrderBook: NewMinHeap(),
		bidsLimitOrderBook: NewMaxHeap(),
		asksStopOrderBook:  NewMaxHeap(),
		bidsStopOrderBook:  NewMinHeap(),
		marketPrice:        decimal.Zero,
		executor:           &Executor{},
	}

	consumer.AddHandler(&MatchHandler{m})
	addr := viper.GetString("nsq.host") + ":" + viper.GetString("nsq.nsqlookupd.port")
	if err := consumer.ConnectToNSQLookupd(addr); err != nil {
		log.Fatalf("[matcher.matcher.NewMatcher] [FETAL] %s", err)
	}

	return m
}

func (m *Matcher) canMatch(order model.Order) bool {
	switch enum.OrdType(order.OrderType) {
	case enum.OrdType_MARKET:
		switch enum.Side(order.Side) {
		case enum.Side_BUY:
			return m.asksLimitOrderBook.Len() != 0
		case enum.Side_SELL:
			return m.bidsLimitOrderBook.Len() != 0
		}
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		return false
	case enum.OrdType_LIMIT:
		switch enum.Side(order.Side) {
		case enum.Side_BUY:
			return m.asksLimitOrderBook.Len() != 0 && order.Price.GreaterThanOrEqual(m.asksLimitOrderBook.Peek().Price)
		case enum.Side_SELL:
			return m.bidsLimitOrderBook.Len() != 0 && order.Price.LessThanOrEqual(m.asksLimitOrderBook.Peek().Price)
		}
	}
	return false
}

// NewMarketOrder creates a new marker order.
func (m *Matcher) NewMarketOrder(order model.Order) {
	var peek *Level
Loop:
	// loop until the open quantity of the order drops to zero
	for order.OpenQuantity.GreaterThan(decimal.Zero) {
		switch enum.Side(order.Side) {
		case enum.Side_BUY:
			if !m.canMatch(order) {
				// asksLimitOrderBook is empty
				if order.Status == string(enum.OrdStatus_PENDING_NEW) {
					// reject invalid orders
					service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_REJECTED))
				}
				break Loop
			}
			peek = m.asksLimitOrderBook.Peek()
		case enum.Side_SELL:
			if !m.canMatch(order) {
				// bidsLimitOrderBook is empty
				if order.Status == string(enum.OrdStatus_PENDING_NEW) {
					// reject invalid orders
					service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_REJECTED))
				}
				break Loop
			}
			peek = m.bidsLimitOrderBook.Peek()
		default:
			log.Print("matcher.matcher.NewMarketOrder [ERROR] Invalid order side.")
			service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_REJECTED))
			break Loop
		}

		if peek.Order[0].OpenQuantity.GreaterThan(order.OpenQuantity) {
			price := peek.Order[0].Price
			quantity := order.OpenQuantity
			// initiator
			peek.Order[0].OpenQuantity = peek.Order[0].OpenQuantity.Sub(quantity)
			peek.Order[0].Status = string(enum.OrdStatus_PARTIALLY_FILLED)
			// completion
			order.OpenQuantity = decimal.Zero
			order.Status = string(enum.OrdStatus_FILLED)

			err := m.executor.NewTrade(peek.Order[0], &order, price, quantity)
			if err != nil {
				prevMarketPrice := m.marketPrice
				m.marketPrice = peek.Order[0].Price
				m.TriggerStopOrder(prevMarketPrice, m.marketPrice)
			}
			break Loop
		}

		price := peek.Order[0].Price
		quantity := peek.Order[0].OpenQuantity
		// completion
		order.OpenQuantity = order.OpenQuantity.Sub(quantity)
		order.Status = string(enum.OrdStatus_PARTIALLY_FILLED)
		// initiator
		peek.Order[0].OpenQuantity = decimal.Zero
		peek.Order[0].Status = string(enum.OrdStatus_FILLED)

		err := m.executor.NewTrade(peek.Order[0], &order, price, quantity)
		if err != nil {
			prevMarketPrice := m.marketPrice
			m.marketPrice = peek.Order[0].Price
			m.TriggerStopOrder(prevMarketPrice, m.marketPrice)
		}

		peek.Order = peek.Order[1:]
		if len(peek.Order) == 0 {
			switch enum.Side(order.Side) {
			case enum.Side_BUY:
				heap.Pop(m.asksLimitOrderBook)
			case enum.Side_SELL:
				heap.Pop(m.bidsLimitOrderBook)
			}
		}
	}
}

// NewLimitOrder creates a new limit order.
func (m *Matcher) NewLimitOrder(order model.Order) {
	switch enum.Side(order.Side) {
	case enum.Side_BUY:
		if !m.canMatch(order) {
			heap.Push(m.bidsLimitOrderBook, Level{order.Price, []*model.Order{&order}})
			service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_NEW))
			break
		}
	case enum.Side_SELL:
		if !m.canMatch(order) {
			heap.Push(m.asksLimitOrderBook, Level{order.Price, []*model.Order{&order}})
			service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_NEW))
			break
		}
	default:
		log.Print("matcher.matcher.NewLimitOrder [ERROR] Invalid order side.")
		service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_REJECTED))
	}
}

// NewStopOrder creates a new stop order.
func (m *Matcher) NewStopOrder(order model.Order) {
	switch enum.Side(order.Side) {
	case enum.Side_BUY:
		if order.StopPrice.GreaterThan(m.marketPrice) {
			heap.Push(m.bidsStopOrderBook, order)
		} else {
			service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_REJECTED))
		}
	case enum.Side_SELL:
		if order.StopPrice.LessThan(m.marketPrice) {
			heap.Push(m.asksStopOrderBook, order)
		} else {
			service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_REJECTED))
		}
	default:
		log.Print("[matcher.matcher.NewStopOrder] [ERROR] Invalid order side.")
	}
}

// NewCancelOrder cancels a specific order.
func (m *Matcher) NewCancelOrder(o model.Order) {
	order := service.Order{}.OrderByID(o.OrderID.String())
	if order.FuturesID == o.FuturesID {
		service.Order{}.CancelOrder(&order)
	}
}

func (m *Matcher) TriggerStopOrder(prev decimal.Decimal, current decimal.Decimal) {
	if prev.GreaterThan(current) {
		for m.asksStopOrderBook.Peek().Price.GreaterThanOrEqual(current) {
			for _, order := range heap.Pop(m.asksStopOrderBook).(Level).Order {
				switch enum.OrdType(order.OrderType) {
				case enum.OrdType_STOP:
					m.NewMarketOrder(*order)
				case enum.OrdType_STOP_LIMIT:
					m.NewLimitOrder(*order)
				}
			}
		}
	} else if prev.LessThan(current) {
		for m.bidsStopOrderBook.Peek().Price.LessThanOrEqual(current) {
			for _, order := range heap.Pop(m.bidsStopOrderBook).(Level).Order {
				switch enum.OrdType(order.OrderType) {
				case enum.OrdType_STOP:
					m.NewMarketOrder(*order)
				case enum.OrdType_STOP_LIMIT:
					m.NewLimitOrder(*order)
				}
			}
		}
	}
}
