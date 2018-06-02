package model

import "github.com/shopspring/decimal"

// Commission represents a service charge assessed by a broker or investment advisor in return for providing investment
// advice and/or handling the purchase or sale of a security.
type Commission struct {
	CommissionID int `gorm:"primary_key;AUTO_INCREMENT"`
	FirmID       int
	FuturesID    string
	OrderType    string
	Percentage   decimal.Decimal `sql:"DECIMAL(10,2)"`
}

func (Commission) TableName() string {
	return "commission"
}
