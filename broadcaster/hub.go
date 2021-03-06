package broadcaster

import (
	"container/heap"
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/matcher"
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/rudeigerc/broker-gateway/service"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type Data interface {
}

type Message struct {
	Type      string `json:"type"`
	Mode      string `json:"mode"`
	FuturesID string `json:"futures_id"`
	Data      Data   `json:"data"`
}

type FuturesData struct {
	Bids [][2]string `json:"bids"`
	Asks [][2]string `json:"asks"`
}

type TradeData struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
	Time     string `json:"time"`
}

type Hub struct {
	EtcdClient *clientv3.Client
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		EtcdClient: mapper.NewEtcdClient(),
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
			tradesMap := service.Trade{}.TradesSnapshot()
			for futuresID, trades := range tradesMap {
				data := make([]TradeData, len(trades))
				for index, trade := range trades {
					data[index] = TradeData{
						Price:    trade.Price.String(),
						Quantity: trade.Quantity.String(),
						Time:     trade.CreatedAt.String(),
					}
				}
				message := Message{
					Type:      "trade",
					Mode:      "snapshot",
					FuturesID: futuresID,
					Data:      data,
				}
				client.message <- message
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.message)
				delete(h.clients, client)
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
	rch := h.EtcdClient.Watch(context.Background(), viper.GetString("etcd.keys.orderbook"), clientv3.WithPrefix(), clientv3.WithProgressNotify())
	for {
		wresp := <-rch

		for _, event := range wresp.Events {
			array := strings.Split(string(event.Kv.Key), "/")
			futuresID := array[len(array)-1]

			asksKey := strings.Replace(viper.GetString("etcd.keys.asks"), "futures_id", futuresID, -1)
			bidsKey := strings.Replace(viper.GetString("etcd.keys.bids"), "futures_id", futuresID, -1)

			asksMarshaled, err := h.EtcdClient.Get(context.Background(), asksKey)
			if err != nil {
				log.Fatalf("[broadcaster.hub.RunOrderBookWatcher] [FETAL] %s", err)
			}
			asksLimitOrderBook := matcher.MinHeap{}
			json.Unmarshal([]byte(asksMarshaled.Kvs[0].Value), &asksLimitOrderBook)
			var asks [][2]string
			for asksLimitOrderBook.Len() > 0 {
				level := heap.Pop(&asksLimitOrderBook).(matcher.Level)
				quantity := decimal.Zero
				for _, order := range level.Order {
					quantity = quantity.Add(order.OpenQuantity)
				}
				asks = append(asks, [2]string{level.Price.String(), quantity.String()})
			}

			bidsMarshaled, err := h.EtcdClient.Get(context.Background(), bidsKey)
			if err != nil {
				log.Fatalf("[broadcaster.hub.RunOrderBookWatcher] [FETAL] %s", err)
			}
			bidsLimitOrderBook := matcher.MaxHeap{}
			json.Unmarshal([]byte(bidsMarshaled.Kvs[0].Value), &bidsLimitOrderBook)
			var bids [][2]string
			for bidsLimitOrderBook.Len() > 0 {
				level := heap.Pop(&bidsLimitOrderBook).(matcher.Level)
				quantity := decimal.Zero
				for _, order := range level.Order {
					quantity = quantity.Add(order.OpenQuantity)
				}
				bids = append(bids, [2]string{level.Price.String(), quantity.String()})
			}

			msg := Message{
				Type:      "order_book",
				Mode:      "snapshot",
				FuturesID: "GC_SEP18",
				Data: FuturesData{
					Bids: bids,
					Asks: asks,
				},
			}
			h.broadcast <- msg
		}
	}
}

func (h *Hub) RunTradeWatcher() {
	etcdClient := mapper.NewEtcdClient()
	defer etcdClient.Close()
	rch := etcdClient.Watch(context.Background(), viper.GetString("etcd.keys.update"), clientv3.WithPrefix())
	for {
		wresp := <-rch

		for _, event := range wresp.Events {
			trade := model.Trade{}
			json.Unmarshal([]byte(event.Kv.Value), &trade)
			msg := Message{
				Type:      "trade",
				Mode:      "update",
				FuturesID: "GC_SEP18",
				Data: TradeData{
					Price:    trade.Price.String(),
					Quantity: trade.Quantity.String(),
					Time:     trade.CreatedAt.Format("2006-01-02 15:04:05 +0800 CST"),
				},
			}
			h.broadcast <- msg
		}
	}
}
