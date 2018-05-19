package model

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

// Order represents an investor's instructions to a broker or brokerage firm to purchase or sell a security.
type Order struct {
	OrderID      uuid.UUID
	OrderType    string
	Side         string
	FuturesID    string
	FirmID       int
	TraderName   string
	Quantity     decimal.Decimal
	OpenQuantity decimal.Decimal
	Price        decimal.Decimal
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
