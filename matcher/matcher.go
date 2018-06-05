package matcher

import (
	"encoding/json"
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
	*Matcher
}

func (h *MatchHandler) HandleMessage(m *nsq.Message) error {
	order := model.Order{}
	json.Unmarshal(m.Body, &order)

	if _, ok := h.MarketDataMap[order.FuturesID]; !ok {
		h.MarketDataMap[order.FuturesID] = NewMarketData(order.FuturesID)
	}

	data := h.MarketDataMap[order.FuturesID]
	switch enum.OrdType(order.OrderType) {
	case enum.OrdType_MARKET:
		service.Order{}.NewOrder(&order)
		data.NewMarketOrder(order)
	case enum.OrdType_LIMIT:
		service.Order{}.NewOrder(&order)
		data.NewLimitOrder(order)
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		service.Order{}.NewOrder(&order)
		data.NewStopOrder(order)
	case enum.OrdType_COUNTER_ORDER_SELECTION:
		data.NewCancelOrder(order)
	default:
		return fmt.Errorf("[matcher.matcher.HandleMessage] [ERROR] Invalid order type: %s", enum.OrdType(order.OrderType))
	}

	return nil
}

type Matcher struct {
	*nsq.Consumer
	MarketDataMap map[string]*MarketData
}

func NewMatcher() *Matcher {
	config := nsq.NewConfig()
	config.LookupdPollInterval = time.Second
	consumer, err := nsq.NewConsumer(viper.GetString("nsq.topic"), "order", config)
	if err != nil {
		log.Fatalf("[matcher.matcher.NewMatcher] [FETAL] %s", err)
	}

	m := &Matcher{
		Consumer:      consumer,
		MarketDataMap: make(map[string]*MarketData),
	}

	consumer.AddHandler(&MatchHandler{m})
	addr := viper.GetString("nsq.host") + ":" + viper.GetString("nsq.nsqlookupd.port")
	if err := consumer.ConnectToNSQLookupd(addr); err != nil {
		log.Fatalf("[matcher.matcher.NewMatcher] [FETAL] %s", err)
	}

	return m
}
