package broadcaster

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/spf13/viper"
)

type Data interface {
}

type Message struct {
	Type      string `json:"type"`
	FuturesID string `json:"futures_id"`
	Data      Data   `json:"data"`
}

type FuturesData struct {
	Bids  [][]string `json:"bids"`
	Asks  [][]string `json:"asks"`
	Level int        `json:"level"`
}

type TradeData struct {
	Price    string    `json:"price"`
	Quantity string    `json:"quantity"`
	Time     time.Time `json:"time"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) RunBroadcaster() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.message)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.message <- message:
				default:
					close(client.message)
					delete(h.clients, client)
				}

			}

		}
	}
}

func (h *Hub) RunOrderBookWatcher() {
	etcdClient := mapper.NewEtcdClient()
	defer etcdClient.Close()

	rch := etcdClient.Watch(context.Background(), "/foo", clientv3.WithPrefix(), clientv3.WithProgressNotify())
	for {
		wresp := <-rch
		fmt.Printf("wresp.Header.Revision: %d\n", wresp.Header.Revision)
		fmt.Println("wresp.IsProgressNotify:", wresp.IsProgressNotify())

		data := FuturesData{
			Bids:  [][]string{{"295.96", "10.34"}},
			Asks:  [][]string{{"295.89", "2.41"}},
			Level: 1,
		}
		msg := Message{
			Type:      "test",
			FuturesID: "test",
			Data:      data,
		}
		h.broadcast <- msg
	}
}

func (h *Hub) RunTradeWatcher() {
	etcdClient := mapper.NewEtcdClient()
	defer etcdClient.Close()

	key := strings.Replace(viper.GetString("etcd.keys.update"), "futures_id", "GC_SEP18", -1)
	rch := etcdClient.Watch(context.Background(), key, clientv3.WithPrefix(), clientv3.WithProgressNotify())
	for {
		<-rch
		marshaled, err := etcdClient.Get(context.Background(), key)
		if err != nil {
			log.Fatalf("[broadcaster.hub.RunTradeWatcher] [FETAL] %s", err)
		}
		trade := model.Trade{}
		json.Unmarshal([]byte(marshaled.Kvs[0].Value), &trade)
		msg := Message{
			Type:      "trade_update",
			FuturesID: "GC_SEP18",
			Data: TradeData{
				Price:    trade.Price.String(),
				Quantity: trade.Quantity.String(),
				Time:     trade.CreatedAt,
			},
		}
		h.broadcast <- msg
	}
}
