package model

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

// Trade represents a basic economic concept involving the buying and selling of goods and services, with compensation
// paid by a buyer to a seller, or the exchange of goods or services between parties.
type Trade struct {
	TradeID              uuid.UUID `gorm:"primary_key"`
	Quantity             decimal.Decimal
	Price                decimal.Decimal
	FuturesID            string
	InitiatorID          int
	InitiatorSide        string
	InitiatorName        string
	InitiatorCommission  decimal.Decimal
	InitiatorOrderID     uuid.UUID
	CompletionID         int
	CompletionSide       string
	CompletionName       string
	CompletionCommission decimal.Decimal
	CompletionOrderID    uuid.UUID
	TradeTime            time.Time
}

func (Trade) TableName() string {
	return "trade"
}
