package broadcaster

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

type Data struct {
	Bids  [][]decimal.Decimal `json:"bids"`
	Asks  [][]decimal.Decimal `json:"asks"`
	Level int                 `json:"level"`
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

func init() {
	go handleBroadcast()
}

func handleBroadcast() {
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
