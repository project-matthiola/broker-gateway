package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rudeigerc/broker-gateway/broadcaster"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func PingHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatalf("[handler.broadcast.PingHandler] [FETAL] %s", err)
	}
	defer ws.Close()

	broadcaster.Clients[ws] = true
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[handler.broadcast.PingHandler] recv: %s", message)

		data := broadcaster.Data{
			Bids:  [][]float64{{295.96, 10.34}},
			Asks:  [][]float64{{295.89, 2.41}},
			Level: 1,
		}
		msg := broadcaster.Message{
			Type:      string(message),
			FuturesId: string(message),
			Data:      data,
		}
		broadcaster.Broadcast <- msg
	}

}
