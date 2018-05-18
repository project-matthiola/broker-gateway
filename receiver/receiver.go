package receiver

import (
	"github.com/quickfixgo/quickfix"
)

type Receiver struct {
	*quickfix.MessageRouter
}

func NewReceiver() *Receiver {
	r := &Receiver{
		MessageRouter: quickfix.NewMessageRouter(),
	}
	return r
}

func (r Receiver) OnCreate(sessionID quickfix.SessionID)                           { return }

func (r Receiver) OnLogon(sessionID quickfix.SessionID)                            { return }

func (r Receiver) OnLogout(sessionID quickfix.SessionID)                           { return }

func (r Receiver) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID)     { return }

func (r Receiver) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error { return nil }

func (r Receiver) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (r Receiver) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return r.Route(msg, sessionID)
}