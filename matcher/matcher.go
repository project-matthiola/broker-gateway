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
	service.Order{}.NewOrder(order)

	switch enum.OrdType(order.OrderType) {
	case enum.OrdType_MARKET:
		h.NewMarketOrder(order)
	case enum.OrdType_LIMIT:
		h.NewLimitOrder(order)
	case enum.OrdType_STOP:
		h.NewStopOrder(order)
	case enum.OrdType_STOP_LIMIT:
		h.NewStopLimitOrder(order)
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

func (m *Matcher) NewMarketOrder(order model.Order) {
	switch enum.Side(order.Side) {
	case enum.Side_BUY:
		if m.canMatch(order) {

		} else {
			service.Order{}.UpdateOrder(order, "status", string(enum.OrdStatus_REJECTED))
		}
	case enum.Side_SELL:
		if m.canMatch(order) {

		} else {
			service.Order{}.UpdateOrder(order, "status", string(enum.OrdStatus_REJECTED))
		}
	default:
		log.Print("matcher.matcher.NewMarketOrder [ERROR] Invalid side of order.")
	}
}

func (m *Matcher) NewLimitOrder(order model.Order) {
	var peek *Level

	for order.OpenQuantity.GreaterThan(decimal.Zero) {
		switch enum.Side(order.Side) {
		case enum.Side_BUY:
			peek = m.asksLimitOrderBook.Peek()
			if !m.canMatch(order) {
				heap.Push(m.bidsLimitOrderBook, Level{order.Price, []*model.Order{&order}})
				service.Order{}.UpdateOrder(order, "status", string(enum.OrdStatus_NEW))
				break
			}
		case enum.Side_SELL:
			peek = m.bidsLimitOrderBook.Peek()
			if !m.canMatch(order) {
				heap.Push(m.asksLimitOrderBook, Level{order.Price, []*model.Order{&order}})
				service.Order{}.UpdateOrder(order, "status", string(enum.OrdStatus_NEW))
				break
			}
		default:
			log.Print("matcher.matcher.NewMarketOrder [ERROR] Invalid side of order.")
			service.Order{}.UpdateOrder(order, "status", string(enum.OrdStatus_REJECTED))
			break
		}

		if peek.Order[0].OpenQuantity.GreaterThan(order.OpenQuantity) {
			peek.Order[0].OpenQuantity = peek.Order[0].OpenQuantity.Sub(order.OpenQuantity)
			break
		}
		order.OpenQuantity = order.OpenQuantity.Sub(peek.Order[0].OpenQuantity)
		peek.Order = peek.Order[1:]
	}
}

func (m *Matcher) NewStopOrder(order model.Order) {
	switch enum.Side(order.Side) {
	case enum.Side_BUY:
		return
	case enum.Side_SELL:
		return
	default:
		log.Print("matcher.matcher.NewMarketOrder [ERROR] Invalid side of order.")
	}
}

func (m *Matcher) NewStopLimitOrder(order model.Order) {
	switch enum.Side(order.Side) {
	case enum.Side_BUY:
		return
	case enum.Side_SELL:
		return
	default:
		log.Print("matcher.matcher.NewMarketOrder [ERROR] Invalid side of order.")
	}
}

func (m *Matcher) NewCancelOrder(order model.Order) {

}
