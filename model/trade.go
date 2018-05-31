package model

import "github.com/pborman/uuid"

// Trade represents a basic economic concept involving the buying and selling of goods and services, with compensation
// paid by a buyer to a seller, or the exchange of goods or services between parties.
type Trade struct {
	TradeId uuid.UUID `gorm:"primary_key"`
}

func (Trade) TableName() string {
	return "trade"
}
