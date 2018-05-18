package model

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

// Order represents an investor's instructions to a broker or brokerage firm to purchase or sell a security.
type Order struct {
	OrderID      uuid.UUID
	OrderType    int
	Side         int
	FuturesID    int
	FirmID       int
	TraderName   string
	Quantity     int
	OpenQuantity int
	Price        decimal.Decimal
	Status       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
