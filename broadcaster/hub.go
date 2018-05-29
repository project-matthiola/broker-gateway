package broadcaster

import (
	"context"
	"fmt"
	"log"

	"github.com/coreos/etcd/clientv3"
	"github.com/gorilla/websocket"
	"github.com/rudeigerc/broker-gateway/mapper"
)

type Data struct {
	Bids  [][]float64 `json:"bids"`
	Asks  [][]float64 `json:"asks"`
	Level int         `json:"level"`
}

type Message struct {
	Type      string `json:"type"`
	FuturesId string `json:"futures_id"`
	Data      Data   `json:"data"`
}

var (
	Clients   = make(map[*websocket.Conn]bool)
	Broadcast = make(chan Message)
)

func HandleBroadcast() {
	for {
		msg := <-Broadcast

		for client := range Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}

func HandleWatcher() {
	etcdClient := mapper.NewEtcdClient()
	rch := etcdClient.Watch(context.Background(), "/foo", clientv3.WithPrefix(), clientv3.WithProgressNotify())
	for {
		wresp := <-rch
		fmt.Printf("wresp.Header.Revision: %d\n", wresp.Header.Revision)
		fmt.Println("wresp.IsProgressNotify:", wresp.IsProgressNotify())

		data := Data{
			Bids:  [][]float64{{295.96, 10.34}},
			Asks:  [][]float64{{295.89, 2.41}},
			Level: 1,
		}
		msg := Message{
			Type:      "test",
			FuturesId: "test",
			Data:      data,
		}
		Broadcast <- msg
	}
}
