package matcher

import (
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/rudeigerc/broker-gateway/service"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type Executor struct {
}

func (e *Executor) NewTrade(initiator model.Order, completion model.Order, price decimal.Decimal, quantity decimal.Decimal) error {
	trade := model.Trade{
		TradeID:              uuid.NewV1(),
		Quantity:             quantity,
		Price:                price,
		FuturesID:            initiator.FuturesID,
		InitiatorID:          initiator.FirmID,
		InitiatorSide:        initiator.Side,
		InitiatorName:        initiator.TraderName,
		InitiatorCommission:  decimal.Zero,
		InitiatorOrderID:     initiator.OrderID,
		CompletionID:         completion.FirmID,
		CompletionSide:       completion.Side,
		CompletionName:       completion.TraderName,
		CompletionCommission: decimal.Zero,
		CompletionOrderID:    completion.OrderID,
	}
	service.Trade{}.NewTrade(&trade)
	service.Order{}.SaveOrder(&initiator)
	service.Order{}.SaveOrder(&completion)
	return nil
}
