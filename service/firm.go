package service

import (
	"log"

	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/model"
)

type Firm struct {
}

func (f Firm) Firms() []model.Firm {
	m := mapper.NewMapper()
	var firms []model.Firm
	err := m.Find(&firms)
	if err != nil {
		log.Printf("[service.firm.Firms] [ERROR] %s", err)
	}
	return firms
}
