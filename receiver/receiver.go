package receiver

import (
	"github.com/quickfixgo/fix50sp2/newordersingle"
	"github.com/quickfixgo/quickfix"
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/satori/go.uuid"
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

	orderType, err := msg.GetOrdType()
	if err != nil {
		return err
	}

	side, err := msg.GetSide()
	if err != nil {
		return err
	}

	futuresID, err := msg.GetProductComplex()
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

	price, err := msg.GetPrice()
	if err != nil {
		return err
	}

	createdAt, err := msg.GetSendingTime()
	if err != nil {
		return err
	}

	order := model.Order{
		OrderID: uuid.NewV1(),
	}

	quickfix.SendToTarget(msg, sessionID)
	return
}
