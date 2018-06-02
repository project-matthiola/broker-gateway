package model

import (
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

// Order represents an investor's instructions to a broker or brokerage firm to purchase or sell a security.
type Order struct {
	OrderID      uuid.UUID `gorm:"primary_key"`
	OrderType    string
	Side         string
	FuturesID    string
	FirmID       int
	TraderName   string
	Quantity     decimal.Decimal `sql:"DECIMAL(10,2)"`
	OpenQuantity decimal.Decimal `sql:"DECIMAL(10,2)"`
	Price        decimal.Decimal `sql:"DECIMAL(10,2)"`
	StopPrice    decimal.Decimal `sql:"DECIMAL(10,2)"`
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Order) TableName() string {
	return "order"
}

func (o Order) Marshal() ([]byte, error) {
	a := []string{
		o.OrderID.String(),
		o.OrderType,
		o.Side,
		o.FuturesID,
		strconv.Itoa(o.FirmID),
		o.TraderName,
		o.Quantity.String(),
		o.OpenQuantity.String(),
		o.Price.String(),
		o.StopPrice.String(),
		o.Status,
		o.CreatedAt.String(),
		o.UpdatedAt.String(),
	}
	s := strings.Join(a, "|")
	return []byte(s), nil
}

func (o *Order) Unmarshal(buf []byte) error {
	var err error

	a := strings.Split(string(buf), "|")

	o.OrderID, err = uuid.FromString(a[0])
	if err != nil {
		return err
	}

	o.OrderType = a[1]

	o.Side = a[2]

	o.FuturesID = a[3]

	o.FirmID, err = strconv.Atoi(a[4])
	if err != nil {
		return err
	}

	o.TraderName = a[5]

	o.Quantity, err = decimal.NewFromString(a[6])
	if err != nil {
		return err
	}

	o.OpenQuantity, err = decimal.NewFromString(a[7])
	if err != nil {
		return err
	}

	o.Price, err = decimal.NewFromString(a[8])
	if err != nil {
		return err
	}

	o.StopPrice, err = decimal.NewFromString(a[9])
	if err != nil {
		return err
	}

	o.Status = a[10]

	o.CreatedAt, err = time.Parse("2006-01-02 15:04:05 +0000 UTC", a[11])
	if err != nil {
		return err
	}

	o.UpdatedAt, err = time.Parse("2006-01-02 15:04:05 +0000 UTC", a[12])
	if err != nil {
		return err
	}

	return nil
}
