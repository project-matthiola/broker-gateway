package cmd

import (
	"log"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/micro/go-micro"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fix50sp2/newordersingle"
	"github.com/quickfixgo/quickfix"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var senderCmd = &cobra.Command{
	Use:   "sender",
	Short: "Run sender",
	Long:  "Run sender",
	Run: func(cmd *cobra.Command, args []string) {
		cfgFileName := path.Join("config", "sender.cfg")

		cfg, err := os.Open(cfgFileName)
		if err != nil {
			log.Printf("[cmd.sender.senderCmd] [ERROR] Error opening %v, %v\n", cfgFileName, err)
			return
		}

		appSettings, err := quickfix.ParseSettings(cfg)
		if err != nil {
			log.Println("[cmd.sender.senderCmd] [ERROR] Error reading cfg,", err)
			return
		}

		app := Sender{}
		fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)

		if err != nil {
			log.Println("[cmd.sender.senderCmd] [ERROR] Error creating file log factory,", err)
			return
		}

		initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, fileLogFactory)
		if err != nil {
			log.Printf("[cmd.sender.senderCmd] [ERROR] Unable to create Initiator: %s\n", err)
			return
		}

		initiator.Start()
		defer initiator.Stop()

		for buy := decimal.NewFromFloat(49.50); buy.LessThanOrEqual(decimal.NewFromFloat(49.99)); buy = buy.Add(decimal.NewFromFloat(0.01)) {
			clOrdID := field.NewClOrdID(uuid.NewV1().String())
			side := field.NewSide(enum.Side_BUY)
			transacttime := field.NewTransactTime(time.Now())
			ordtype := field.NewOrdType(enum.OrdType_LIMIT)

			order := newordersingle.New(clOrdID, side, transacttime, ordtype)
			order.SetSenderCompID("Sender")
			order.SetSenderSubID("Rudeiger Cheng")
			order.SetTargetCompID("Broker")
			order.SetSymbol("GC_SEP18")
			order.SetPrice(buy, 2)
			order.SetOrderQty(decimal.NewFromFloat(1000).Div(buy).Mul(decimal.NewFromFloat(rand.Float64())), 2)
			msg := order.ToMessage()
			quickfix.Send(msg)
		}

		for sell := decimal.NewFromFloat(50.41); sell.GreaterThanOrEqual(decimal.NewFromFloat(50.00)); sell = sell.Sub(decimal.NewFromFloat(0.01)) {
			clOrdID := field.NewClOrdID(uuid.NewV1().String())
			side := field.NewSide(enum.Side_SELL)
			transacttime := field.NewTransactTime(time.Now())
			ordtype := field.NewOrdType(enum.OrdType_LIMIT)

			order := newordersingle.New(clOrdID, side, transacttime, ordtype)
			order.SetSenderCompID("Sender")
			order.SetSenderSubID("Rudeiger Cheng")
			order.SetTargetCompID("Broker")
			order.SetSymbol("GC_SEP18")
			order.SetPrice(sell, 2)
			order.SetOrderQty(decimal.NewFromFloat(1000).Div(sell).Mul(decimal.NewFromFloat(rand.Float64())), 2)
			msg := order.ToMessage()
			quickfix.Send(msg)
		}

		service := micro.NewService(
			micro.Name("github.com.rudeigerc.broker-gateway.sender"),
			micro.RegisterTTL(time.Minute),
			micro.RegisterInterval(time.Second*30),
		)

		if err := service.Run(); err != nil {
			log.Fatalf("[cmd.sender.senderCmd] [FETAL] %s", err)
		}
	},
}

type Sender struct {
}

// OnCreate implemented as part of Application interface.
func (r Sender) OnCreate(sessionID quickfix.SessionID) { return }

// OnLogon implemented as part of Application interface.
func (r Sender) OnLogon(sessionID quickfix.SessionID) { return }

// OnLogout implemented as part of Application interface.
func (r Sender) OnLogout(sessionID quickfix.SessionID) { return }

// ToAdmin implemented as part of Application interface.
func (r Sender) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) { return }

// ToApp implemented as part of Application interface.
func (r Sender) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	log.Printf("[cmd.sender.ToApp] Sending %s\n", msg)
	return nil
}

// FromAdmin implemented as part of Application interface
func (r Sender) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

// FromApp implemented as part of Application interface.
func (r Sender) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	log.Printf("[cmd.sender.FromApp] FromApp: %s\n", msg.String())
	return nil
}
