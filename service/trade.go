package service

import (
	"log"

	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/model"
)

type Trade struct {
}

func (t Trade) NewTrade(trade *model.Trade) {
	m := mapper.NewMapper()
	err := m.Create(trade)
	if err != nil {
		log.Printf("[service.order.NewOrder] [ERROR] %s", err)
	}
}

func (t Trade) Trades() []model.Trade {
	m := mapper.NewMapper()
	var trades []model.Trade
	err := m.FindWithLimit(&trades, -1)
	if err != nil {
		log.Printf("[service.trade.Trades] [ERROR] %s", err)
	}
	return trades
}

func (t Trade) TradeByID(uuid string) model.Trade {
	m := mapper.NewMapper()
	trade := model.Trade{}
	err := m.WhereByUUID(&trade, "trade_id", uuid)
	if err != nil {
		log.Printf("[service.order.OrderByID] [ERROR] %s", err)
	}
	return trade
}

func (t Trade) TradesSnapshot() map[string][]model.Trade {
	m := mapper.NewMapper()
	futuresIDs, err := m.FutureIDs()
	if err != nil {
		log.Printf("[service.trade.TradesSnapshot] [ERROR] %s", err)
	}

	tradesMap := make(map[string][]model.Trade)
	for _, futuresID := range futuresIDs {
		var trades []model.Trade
		err := m.FindByFuturesID(&trades, futuresID)
		if err != nil {
			log.Printf("[service.trade.TradesSnapshot] [ERROR] %s", err)
		}
		tradesMap[futuresID] = trades
	}

	return tradesMap
}

func (t Trade) TradesWithPage(page int) []model.Trade {
	m := mapper.NewMapper()
	var trades []model.Trade
	err := m.FindWithPage(&trades, page)
	if err != nil {
		log.Printf("[service.trade.OrdersWithPage] [ERROR] %s", err)
	}
	return trades
}

func (t Trade) TradesWithCondition(firmID int, futuresID string, traderName string, page int) []model.Trade {
	m := mapper.NewMapper()
	var trades []model.Trade
	err := m.FindTradesWithCondition(&trades, firmID, futuresID, traderName, page)
	if err != nil {
		log.Printf("[service.trade.TradesWithCondition] [ERROR] %s", err)
	}
	return trades
}
