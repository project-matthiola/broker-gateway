package model

import (
	"time"

	"strconv"

	"strings"

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

	o.Status = a[9]

	o.CreatedAt, err = time.Parse("2006-01-02 15:04:05.000 +0000 UTC", a[10])
	if err != nil {
		return err
	}

	o.UpdatedAt, err = time.Parse("2006-01-02 15:04:05.000 +0000 UTC", a[11])
	if err != nil {
		return err
	}

	return nil
}
