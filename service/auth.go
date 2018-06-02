package service

import (
	"log"

	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/model"
)

type Auth struct {
}

func (a Auth) Sign(firmName string) model.Firm {
	m := mapper.NewMapper()
	firm := model.Firm{
		FirmName: firmName,
	}
	err := m.Create(&firm)
	if err != nil {
		log.Printf("[service.auth] [ERROR] %s", err)
	}
	return firm
}

func (a Auth) Validate(firmID int) model.Firm {
	m := mapper.NewMapper()
	firm := model.Firm{}
	err := m.FirstByID(&firm, firmID)
	if err != nil {
		log.Printf("[service.auth] [ERROR] %s", err)
	}
	return firm
}
