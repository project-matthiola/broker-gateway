package service

import (
	"log"

	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/model"
)

type Order struct {
}

func (o Order) NewOrder(order model.Order) {
	m := mapper.NewMapper()
	err := m.Create(&order)
	if err != nil {
		log.Fatal(err)
	}
}
