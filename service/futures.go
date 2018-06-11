package service

import (
	"log"

	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/model"
)

type Futures struct {
}

func (f Futures) Futures() map[string][]model.Futures {
	m := mapper.NewMapper()
	futuresNames, err := m.FuturesNames()
	if err != nil {
		log.Printf("[service.futures.Futures] [ERROR] %s", err)
	}

	futuresMap := make(map[string][]model.Futures)
	for _, futuresName := range futuresNames {
		var futures []model.Futures
		err := m.FindByFuturesName(&futures, futuresName)
		if err != nil {
			log.Printf("[service.futures.Futures] [ERROR] %s", err)
		}
		futuresMap[futuresName] = futures
	}

	return futuresMap
}
