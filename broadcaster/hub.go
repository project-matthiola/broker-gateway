package broadcaster

import (
	"log"

	"github.com/gorilla/websocket"
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
