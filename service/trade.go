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
