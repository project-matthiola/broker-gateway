package broadcaster

import (
	"context"
	"fmt"

	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/rudeigerc/broker-gateway/mapper"
)

type Message interface {
}

type Data struct {
	Bids  [][]string `json:"bids"`
	Asks  [][]string `json:"asks"`
	Level int        `json:"level"`
}

type Futures struct {
	Type      string `json:"type"`
	FuturesID string `json:"futures_id"`
	Data      Data   `json:"data"`
}

type TradeData struct {
	Price    string    `json:"price"`
	Quantity string    `json:"quantity"`
	Time     time.Time `json:"time"`
}

type Trade struct {
	Type      string `json:"type"`
	FuturesID string `json:"futures_id"`
	Data      []Data `json:"data"`
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

func (h *Hub) RunWatcher() {
	etcdClient := mapper.NewEtcdClient()
	defer etcdClient.Close()

	rch := etcdClient.Watch(context.Background(), "/foo", clientv3.WithPrefix(), clientv3.WithProgressNotify())
	for {
		wresp := <-rch
		fmt.Printf("wresp.Header.Revision: %d\n", wresp.Header.Revision)
		fmt.Println("wresp.IsProgressNotify:", wresp.IsProgressNotify())

		data := Data{
			Bids:  [][]string{{"295.96", "10.34"}},
			Asks:  [][]string{{"295.89", "2.41"}},
			Level: 1,
		}
		msg := Futures{
			Type:      "test",
			FuturesID: "test",
			Data:      data,
		}
		h.broadcast <- msg
	}
}
