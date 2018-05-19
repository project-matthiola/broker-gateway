package receiver

import (
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/fix50sp2/newordersingle"
	"github.com/quickfixgo/quickfix"
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/satori/go.uuid"

	"log"
	"strconv"
)

type Receiver struct {
	*quickfix.MessageRouter
}

func NewReceiver() *Receiver {
	r := &Receiver{
		MessageRouter: quickfix.NewMessageRouter(),
	}
	r.AddRoute(newordersingle.Route(r.OnNewOrderSingle))
	return r
}

func (r Receiver) OnCreate(sessionID quickfix.SessionID) { return }

func (r Receiver) OnLogon(sessionID quickfix.SessionID) { return }

func (r Receiver) OnLogout(sessionID quickfix.SessionID) { return }

func (r Receiver) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) { return }

func (r Receiver) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error { return nil }

func (r Receiver) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (r *Receiver) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return r.Route(msg, sessionID)
}

func (r *Receiver) OnNewOrderSingle(msg newordersingle.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {

	ordType, err := msg.GetOrdType()
	if err != nil {
		return err
	}

	side, err := msg.GetSide()
	if err != nil {
		return err
	}

	futuresID, err := msg.GetSymbol()
	if err != nil {
		return err
	}

	firmID, err := msg.GetSenderCompID()
	if err != nil {
		return err
	}

	traderName, err := msg.GetSenderSubID()
	if err != nil {
		return err
	}

	quantity, err := msg.GetOrderQty()
	if err != nil {
		return err
	}

	price, _ := msg.GetPrice()

	createdAt, err := msg.GetTransactTime()
	if err != nil {
		return err
	}

	firmIDInt, _ := strconv.Atoi(firmID)

	order := model.Order{
		OrderID:      uuid.NewV1(),
		OrderType:    string(ordType),
		Side:         string(side),
		FuturesID:    futuresID,
		FirmID:       firmIDInt,
		TraderName:   traderName,
		Quantity:     quantity,
		OpenQuantity: quantity,
		Price:        price,
		Status:       string(enum.OrdStatus_NEW),
		CreatedAt:    createdAt,
		UpdatedAt:    createdAt,
	}

	log.Print(order)

	quickfix.SendToTarget(msg, sessionID)
	return
}
