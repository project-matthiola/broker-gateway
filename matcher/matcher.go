package matcher

import (
	"fmt"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/quickfixgo/enum"
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/rudeigerc/broker-gateway/service"
	"github.com/spf13/viper"
)

type MatchHandler struct {
}

func (h *MatchHandler) HandleMessage(m *nsq.Message) error {
	order := model.Order{}
	order.Unmarshal(m.Body)
	service.Order{}.NewOrder(order)

	switch enum.OrdType(order.OrderType) {
	case enum.OrdType_MARKET:
		NewMarketOrder(order)
	case enum.OrdType_LIMIT:
		NewLimitOrder(order)
	case enum.OrdType_STOP:
		NewStopOrder(order)
	case enum.OrdType_COUNTER_ORDER_SELECTION:
		NewCancelOrder(order)
	default:
		return fmt.Errorf("[matcher.matcher.HandleMessage] [ERROR] Invalid order type: %s", enum.OrdType(order.OrderType))
	}

	return nil
}

type Matcher struct {
	*nsq.Consumer
}

func NewMatcher() *Matcher {
	config := nsq.NewConfig()
	config.LookupdPollInterval = time.Second
	consumer, err := nsq.NewConsumer(viper.GetString("nsq.topic"), "order", config)
	if err != nil {
		log.Fatalf("[matcher.matcher.NewMatcher] [FETAL] %s", err)
	}

	m := &Matcher{
		Consumer: consumer,
	}

	consumer.AddHandler(&MatchHandler{})
	addr := viper.GetString("nsq.host") + ":" + viper.GetString("nsq.nsqlookupd.port")
	if err := consumer.ConnectToNSQLookupd(addr); err != nil {
		log.Fatalf("[matcher.matcher.NewMatcher] [FETAL] %s", err)
	}

	return m
}

func NewMarketOrder(order model.Order) {
	switch enum.Side(order.Side) {
	case enum.Side_BUY:
		return
	case enum.Side_SELL:
		return
	default:
		log.Print("matcher.matcher.NewMarketOrder [ERROR] Invalid side of order.")
	}
}

func NewLimitOrder(order model.Order) {
	switch enum.Side(order.Side) {
	case enum.Side_BUY:
		return
	case enum.Side_SELL:
		return
	default:
		log.Print("matcher.matcher.NewMarketOrder [ERROR] Invalid side of order.")
	}
}

func NewStopOrder(order model.Order) {
	switch enum.Side(order.Side) {
	case enum.Side_BUY:
		return
	case enum.Side_SELL:
		return
	default:
		log.Print("matcher.matcher.NewMarketOrder [ERROR] Invalid side of order.")
	}
}

func NewCancelOrder(order model.Order) {

}
