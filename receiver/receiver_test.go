package receiver

import (
	"testing"

	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fix50sp2/newordersingle"
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

	msg := newordersingle.New(clOrdID, side, transacttime, ordtype)
	msg.SetSymbol("GC_SEP18")
	msg.SetSenderCompID("1")
	msg.SetApplVerID(enum.ApplVerID_FIX50SP2)
	msg.SetSenderSubID("John Doe")
	msg.SetOrderQty(decimal.NewFromFloat(23.14), 2)

	session := quickfix.SessionID{}
	err := receiver.OnNewOrderSingle(msg, session)
	if err != nil {
		t.Error(err)
	}
}
