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
		log.Printf("[service.auth.Sign] [ERROR] %s", err)
	}
	return firm
}

func (a Auth) Validate(firmID int) model.Firm {
	m := mapper.NewMapper()
	firm := model.Firm{}
	err := m.FirstByID(&firm, firmID)
	if err != nil {
		log.Printf("[service.auth.Validate] [ERROR] %s", err)
	}
	return firm
}

func (a Auth) FirmNameByID(firmID int) string {
	return a.Validate(firmID).FirmName
}

func (a Auth) FirmIDByName(firmName string) int {
	m := mapper.NewMapper()
	firm := model.Firm{}
	err := m.FirstByFirmName(&firm, firmName)
	if err != nil {
		log.Printf("[service.auth.FirmIDByName] [ERROR] %s", err)
	}
	return firm.FirmID
}
