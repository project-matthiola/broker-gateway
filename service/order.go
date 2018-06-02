package service

import (
	"log"

	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/model"
)

type Order struct {
}

func (o Order) NewOrder(order *model.Order) {
	m := mapper.NewMapper()
	err := m.Create(order)
	if err != nil {
		log.Fatalf("[service.order.NewOrder] [FETAL] %s", err)
	}
}

func (o Order) SaveOrder(order *model.Order) {
	m := mapper.NewMapper()
	err := m.Save(order)
	if err != nil {
		log.Fatalf("[service.order.SaveOrder] [FETAL] %s", err)
	}
}

func (o Order) UpdateOrder(order *model.Order, column string, value string) {
	m := mapper.NewMapper()
	err := m.Update(order, column, value)
	if err != nil {
		log.Fatalf("[service.order.UpdateOrder] [FETAL] %s", err)
	}
}
