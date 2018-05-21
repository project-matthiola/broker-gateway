package receiver

import (
	"log"
	"strconv"

	"github.com/nsqio/go-nsq"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fix50sp2/executionreport"
	"github.com/quickfixgo/fix50sp2/newordersingle"
	"github.com/quickfixgo/fix50sp2/ordercancelrequest"
	"github.com/quickfixgo/quickfix"
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

// Receiver implements the quickfix.Application interface.
type Receiver struct {
	*quickfix.MessageRouter
	*nsq.Producer
}

// NewReceiver returns a new receiver.
func NewReceiver() *Receiver {
	nsqAddr := viper.GetString("nsq.host") + ":" + viper.GetString("nsq.port")
	producer, err := nsq.NewProducer(nsqAddr, nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	r := &Receiver{
		MessageRouter: quickfix.NewMessageRouter(),
		Producer:      producer,
	}
	r.AddRoute(newordersingle.Route(r.OnNewOrderSingle))
	r.AddRoute(ordercancelrequest.Route(r.OnOrderCancelRequest))
	return r
}

// OnCreate implemented as part of Application interface.
func (r Receiver) OnCreate(sessionID quickfix.SessionID) { return }

// OnLogon implemented as part of Application interface.
func (r Receiver) OnLogon(sessionID quickfix.SessionID) { return }

// OnLogout implemented as part of Application interface.
func (r Receiver) OnLogout(sessionID quickfix.SessionID) { return }

// ToAdmin implemented as part of Application interface.
func (r Receiver) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) { return }

// ToApp implemented as part of Application interface.
func (r Receiver) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error { return nil }

// FromAdmin implemented as part of Application interface
func (r Receiver) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

// FromApp implemented as part of Application interface, uses Router on incoming application messages.
func (r *Receiver) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return r.Route(msg, sessionID)
}

// OnNewOrderSingle handles the NewOrderSingle.
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
	orderID := uuid.NewV1()

	order := model.Order{
		OrderID:      orderID,
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

	marshaled, _ := order.Marshal()
	r.Publish(viper.GetString("nsq.topic"), marshaled)

	execReport := executionreport.New(
		field.NewOrderID(order.OrderID.String()),
		field.NewExecID(uuid.NewV1().String()),
		field.NewExecType(enum.ExecType(enum.OrdStatus_NEW)),
		field.NewOrdStatus(enum.OrdStatus_NEW),
		field.NewSide(side),
		field.NewLeavesQty(order.OpenQuantity, 2),
		field.NewCumQty(decimal.Zero, 2),
	)

	execReport.SetOrderQty(order.Quantity, 2)

	quickfix.SendToTarget(execReport, sessionID)
	return nil
}

func (r *Receiver) OnOrderCancelRequest(msg ordercancelrequest.OrderCancelRequest, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	orderID, err := msg.GetOrderID()
	if err != nil {
		return err
	}

	symbol, err := msg.GetSymbol()
	if err != nil {
		return err
	}

	orderUUID, _ := uuid.FromString(orderID)

	order := model.Order{
		OrderID:   orderUUID,
		OrderType: string(enum.OrdType_COUNTER_ORDER_SELECTION),
		FuturesID: symbol,
	}

	log.Print(order)

	return nil
}
