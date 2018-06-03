package service

import (
	"log"

	"github.com/quickfixgo/enum"
	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/model"
)

type Order struct {
}

func (o Order) NewOrder(order *model.Order) {
	m := mapper.NewMapper()
	err := m.Create(order)
	if err != nil {
		log.Printf("[service.order.NewOrder] [ERROR] %s", err)
	}
}

func (o Order) SaveOrder(order *model.Order) {
	m := mapper.NewMapper()
	err := m.Save(order)
	if err != nil {
		log.Printf("[service.order.SaveOrder] [ERROR] %s", err)
	}
}

func (o Order) UpdateOrder(order *model.Order, column string, value string) {
	m := mapper.NewMapper()
	err := m.Update(order, column, value)
	if err != nil {
		log.Printf("[service.order.UpdateOrder] [ERROR] %s", err)
	}
}

func (o Order) CancelOrder(order *model.Order) {
	m := mapper.NewMapper()
	err := m.Update(order, "status", string(enum.OrdStatus_CANCELED))
	err = m.Delete(order)
	if err != nil {
		log.Printf("[service.order.CancelOrder] [ERROR] %s", err)
	}
}

func (o Order) OrderByID(uuid string) model.Order {
	m := mapper.NewMapper()
	order := model.Order{}
	err := m.WhereByUUID(&order, "order_id", uuid)
	if err != nil {
		log.Printf("[service.order.OrderByID] [ERROR] %s", err)
	}
	return order
}

func (o Order) Orders() []model.Order {
	m := mapper.NewMapper()
	var orders []model.Order
	err := m.FindWithLimit(&orders, -1)
	if err != nil {
		log.Printf("[service.order.Orders] [ERROR] %s", err)
	}
	return orders
}
