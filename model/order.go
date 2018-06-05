package model

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

// Order represents an investor's instructions to a broker or brokerage firm to purchase or sell a security.
type Order struct {
	OrderID      uuid.UUID       `json:"order_id" gorm:"primary_key"`
	OrderType    string          `json:"order_type"`
	Side         string          `json:"side"`
	FuturesID    string          `json:"futures_id"`
	FirmID       int             `json:"firm_id"`
	TraderName   string          `json:"trader_name"`
	Quantity     decimal.Decimal `json:"quantity" sql:"DECIMAL(10,2)"`
	OpenQuantity decimal.Decimal `json:"open_quantity" sql:"DECIMAL(10,2)"`
	Price        decimal.Decimal `json:"price" sql:"DECIMAL(10,2)"`
	StopPrice    decimal.Decimal `json:"stop_price" sql:"DECIMAL(10,2)"`
	Status       string          `json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func (Order) TableName() string {
	return "order"
}
