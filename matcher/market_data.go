package matcher

import (
	"container/heap"
	"context"
	"log"
	"strings"

	"github.com/gin-gonic/gin/json"
	"github.com/quickfixgo/enum"
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/rudeigerc/broker-gateway/service"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type MarketData struct {
	futuresID string

	asksLimitOrderBook *MinHeap
	bidsLimitOrderBook *MaxHeap

	asksStopOrderBook *MaxHeap
	bidsStopOrderBook *MinHeap

	marketPrice decimal.Decimal

	Executor *Executor
}

func NewMarketData(futuresID string) *MarketData {
	return &MarketData{
		futuresID:          futuresID,
		asksLimitOrderBook: NewMinHeap(),
		bidsLimitOrderBook: NewMaxHeap(),
		asksStopOrderBook:  NewMaxHeap(),
		bidsStopOrderBook:  NewMinHeap(),
		marketPrice:        decimal.Zero,
		Executor:           NewExecutor(),
	}
}

func (m *MarketData) canMatch(order model.Order) bool {
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
			return m.bidsLimitOrderBook.Len() != 0 && order.Price.LessThanOrEqual(m.bidsLimitOrderBook.Peek().Price)
		}
	}
	return false
}

// NewMarketOrder creates a new marker order.
func (m *MarketData) NewMarketOrder(order model.Order) {
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

			err := m.Executor.NewTrade(peek.Order[0], &order, price, quantity)
			if err == nil {
				prevMarketPrice := m.marketPrice
				m.marketPrice = peek.Order[0].Price
				m.BroadcastOrderBook()
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

		err := m.Executor.NewTrade(peek.Order[0], &order, price, quantity)
		if err == nil {
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
		m.BroadcastOrderBook()
	}
}

// NewLimitOrder creates a new limit order.
func (m *MarketData) NewLimitOrder(order model.Order) {
	var peek *Level
	for order.OpenQuantity.GreaterThan(decimal.Zero) {
		switch enum.Side(order.Side) {
		case enum.Side_BUY:
			if !m.canMatch(order) {
				heap.Push(m.bidsLimitOrderBook, Level{order.Price, []*model.Order{&order}})
				service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_NEW))
				m.BroadcastOrderBook()
			}
			return
		case enum.Side_SELL:
			if !m.canMatch(order) {
				heap.Push(m.asksLimitOrderBook, Level{order.Price, []*model.Order{&order}})
				service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_NEW))
				m.BroadcastOrderBook()
			}
			return
		default:
			log.Print("matcher.matcher.NewLimitOrder [ERROR] Invalid order side.")
			service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_REJECTED))
			return
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

			err := m.Executor.NewTrade(peek.Order[0], &order, price, quantity)
			if err == nil {
				prevMarketPrice := m.marketPrice
				m.marketPrice = peek.Order[0].Price
				m.BroadcastOrderBook()
				m.TriggerStopOrder(prevMarketPrice, m.marketPrice)
			}
			return
		}

		price := peek.Order[0].Price
		quantity := peek.Order[0].OpenQuantity
		// completion
		order.OpenQuantity = order.OpenQuantity.Sub(quantity)
		order.Status = string(enum.OrdStatus_PARTIALLY_FILLED)
		// initiator
		peek.Order[0].OpenQuantity = decimal.Zero
		peek.Order[0].Status = string(enum.OrdStatus_FILLED)

		err := m.Executor.NewTrade(peek.Order[0], &order, price, quantity)
		if err == nil {
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
		m.BroadcastOrderBook()
	}
}

// NewStopOrder creates a new stop order.
func (m *MarketData) NewStopOrder(order model.Order) {
	switch enum.Side(order.Side) {
	case enum.Side_BUY:
		if order.StopPrice.GreaterThan(m.marketPrice) {
			heap.Push(m.bidsStopOrderBook, Level{order.Price, []*model.Order{&order}})
		} else {
			service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_REJECTED))
		}
	case enum.Side_SELL:
		if order.StopPrice.LessThan(m.marketPrice) {
			heap.Push(m.asksStopOrderBook, Level{order.Price, []*model.Order{&order}})
		} else {
			service.Order{}.UpdateOrder(&order, "status", string(enum.OrdStatus_REJECTED))
		}
	default:
		log.Print("[matcher.matcher.NewStopOrder] [ERROR] Invalid order side.")
	}
}

// NewCancelOrder cancels a specific order.
func (m *MarketData) NewCancelOrder(o model.Order) {
	var err error
	order := service.Order{}.OrderByID(o.OrderID.String())

	if order.FuturesID == o.FuturesID {
		switch enum.OrdType(order.OrderType) {
		case enum.OrdType_STOP:
			switch enum.Side(order.Side) {
			case enum.Side_BUY:
				err = m.bidsStopOrderBook.Remove(order)
			case enum.Side_SELL:
				err = m.asksStopOrderBook.Remove(order)
			}
		case enum.OrdType_STOP_LIMIT:
			switch enum.Side(order.Side) {
			case enum.Side_BUY:
				err = m.bidsStopOrderBook.Remove(order)
			case enum.Side_SELL:
				err = m.asksStopOrderBook.Remove(order)
			}
			if err != nil {
				switch enum.Side(order.Side) {
				case enum.Side_BUY:
					err = m.bidsLimitOrderBook.Remove(order)
				case enum.Side_SELL:
					err = m.asksLimitOrderBook.Remove(order)
				}
			}
		case enum.OrdType_LIMIT:
			switch enum.Side(order.Side) {
			case enum.Side_BUY:
				err = m.bidsLimitOrderBook.Remove(order)
			case enum.Side_SELL:
				err = m.asksLimitOrderBook.Remove(order)
			}
		}
		service.Order{}.CancelOrder(&order)
	}
	m.BroadcastOrderBook()
}

func (m *MarketData) TriggerStopOrder(prev decimal.Decimal, current decimal.Decimal) {
	if prev.GreaterThan(current) && m.asksStopOrderBook.Len() != 0 {
		for m.asksStopOrderBook.Peek().Price.GreaterThanOrEqual(current) {
			for _, order := range heap.Pop(m.asksStopOrderBook).(Level).Order {
				switch enum.OrdType(order.OrderType) {
				case enum.OrdType_STOP:
					m.NewMarketOrder(*order)
				case enum.OrdType_STOP_LIMIT:
					m.NewLimitOrder(*order)
				}
				if m.asksStopOrderBook.Len() == 0 {
					return
				}
			}
		}
	} else if prev.LessThan(current) && m.bidsStopOrderBook.Len() != 0 {
		for m.bidsStopOrderBook.Peek().Price.LessThanOrEqual(current) {
			for _, order := range heap.Pop(m.bidsStopOrderBook).(Level).Order {
				switch enum.OrdType(order.OrderType) {
				case enum.OrdType_STOP:
					m.NewMarketOrder(*order)
				case enum.OrdType_STOP_LIMIT:
					m.NewLimitOrder(*order)
				}
				if m.bidsStopOrderBook.Len() == 0 {
					return
				}
			}
		}
	}
}

func (m *MarketData) BroadcastOrderBook() {
	asksMarshaled, err := json.Marshal(m.asksLimitOrderBook)
	if err != nil {
		log.Printf("[matcher.market_data.BroadcastOrderBook] [ERROR] %s", err)
	}

	bidsMarshaled, err := json.Marshal(m.bidsLimitOrderBook)
	if err != nil {
		log.Printf("[matcher.market_data.BroadcastOrderBook] [ERROR] %s", err)
	}

	asksKey := strings.Replace(viper.GetString("etcd.keys.asks"), "futures_id", m.futuresID, -1)
	bidsKey := strings.Replace(viper.GetString("etcd.keys.bids"), "futures_id", m.futuresID, -1)
	m.Executor.EtcdClient.Put(context.Background(), asksKey, string(asksMarshaled))
	m.Executor.EtcdClient.Put(context.Background(), bidsKey, string(bidsMarshaled))
}
