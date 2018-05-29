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
	// log.Println(order)
	service.Order{}.NewOrder(order)

	switch enum.OrdType(order.OrderType) {
	case enum.OrdType_MARKET:
		newMarketOrder(order)
	case enum.OrdType_LIMIT:
		newLimitOrder(order)
	case enum.OrdType_STOP:
		newStopOrder(order)
	case enum.OrdType_COUNTER_ORDER_SELECTION:
		newCancelOrder(order)
	default:
		return fmt.Errorf("[matcher] [ERROR] Invalid order type: %s", enum.OrdType(order.OrderType))
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
		log.Fatal(err)
	}

	m := &Matcher{
		Consumer: consumer,
	}

	consumer.AddHandler(&MatchHandler{})
	addr := viper.GetString("nsq.host") + ":" + viper.GetString("nsq.nsqlookupd.port")
	if err := consumer.ConnectToNSQLookupd(addr); err != nil {
		panic(err)
	}

	return m
}

func newMarketOrder(order model.Order) {

}

func newLimitOrder(order model.Order) {

}

func newStopOrder(order model.Order) {

}

func newCancelOrder(order model.Order) {

}
