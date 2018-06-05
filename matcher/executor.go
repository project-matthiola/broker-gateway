package matcher

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/rudeigerc/broker-gateway/service"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type Executor struct {
	EtcdClient *clientv3.Client
}

func NewExecutor() *Executor {
	return &Executor{
		EtcdClient: mapper.NewEtcdClient(),
	}
}

func (e *Executor) NewTrade(initiator *model.Order, completion *model.Order, price decimal.Decimal, quantity decimal.Decimal) error {
	trade := &model.Trade{
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
	service.Trade{}.NewTrade(trade)
	marshaled, err := json.Marshal(trade)
	if err != nil {
		log.Panicf("[matcher.executor.NewTrade] [ERROR} %s", err)
	}

	key := strings.Replace(viper.GetString("etcd.keys.update"), "futures_id", trade.FuturesID, -1)
	e.EtcdClient.Put(context.Background(), key, string(marshaled))
	service.Order{}.SaveOrder(initiator)
	service.Order{}.SaveOrder(completion)
	return nil
}
