package model

import "github.com/shopspring/decimal"

// Commission represents a service charge assessed by a broker or investment advisor in return for providing investment
// advice and/or handling the purchase or sale of a security.
type Commission struct {
	CommissionID int `gorm:"primary_key"`
	FirmID       int
	FuturesID    string
	OrderType    string
	Percentage   decimal.Decimal
}

func (Commission) TableName() string {
	return "commission"
}
