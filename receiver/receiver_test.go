package receiver

import (
	"testing"
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fix50sp2/newordersingle"
	"github.com/quickfixgo/fix50sp2/ordercancelrequest"
	"github.com/quickfixgo/quickfix"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

var receiver *Receiver

func TestNewReceiver(t *testing.T) {
	receiver = NewReceiver()
}

func TestReceiver_OnNewOrderSingle(t *testing.T) {

	clOrdID := field.NewClOrdID(uuid.NewV1().String())
	side := field.NewSide(enum.Side_BUY)
	transacttime := field.NewTransactTime(time.Now())
	ordtype := field.NewOrdType(enum.OrdType_MARKET)

	order := newordersingle.New(clOrdID, side, transacttime, ordtype)
	order.SetApplVerID(enum.ApplVerID_FIX50SP2)
	order.SetSenderCompID("Trader")
	order.SetSenderSubID("John Doe")
	order.SetTargetCompID("Broker")
	order.SetSymbol("GC_SEP18")
	order.SetOrderQty(decimal.NewFromFloat(23.14), 2)

	session := quickfix.SessionID{}
	err := receiver.OnNewOrderSingle(order, session)
	if err != nil {
		t.Error(err)
	}
}

func TestReceiver_OnOrderCancelRequest(t *testing.T) {

	clOrdID := field.NewClOrdID(uuid.NewV1().String())
	side := field.NewSide(enum.Side_BUY)
	transacttime := field.NewTransactTime(time.Now())

	order := ordercancelrequest.New(clOrdID, side, transacttime)
	order.SetApplVerID(enum.ApplVerID_FIX50SP2)
	order.SetSenderCompID("Trader")
	order.SetSenderSubID("John Doe")
	order.SetTargetCompID("Broker")
	order.SetOrderID(uuid.NewV1().String())
	order.SetSymbol("GC_SEP18")

	session := quickfix.SessionID{}
	err := receiver.OnOrderCancelRequest(order, session)
	if err != nil {
		t.Error(err)
	}
}
