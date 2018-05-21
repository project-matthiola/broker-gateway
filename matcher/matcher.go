package matcher

import (
	"log"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/spf13/viper"
)

type MatchHandler struct {
}

func (h *MatchHandler) HandleMessage(m *nsq.Message) error {
	order := model.Order{}
	order.Unmarshal(m.Body)
	log.Println(order)
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
