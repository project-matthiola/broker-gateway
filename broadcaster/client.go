package broadcaster

import (
	"log"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	futuresID string
	hub       *Hub
	conn      *websocket.Conn
	message   chan Message
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SocketHandler(hub *Hub, c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatalf("[broadcast.client.SocketHandler] [FETAL] %s", err)
	}

	client := &Client{strings.Trim(c.Param("futures_id"), "/"), hub, ws, make(chan Message)}
	client.hub.register <- client

	go client.writeMessage()
	go client.readMessage()
}

func (c *Client) writeMessage() {
	defer c.conn.Close()
	for {
		select {
		case message, ok := <-c.message:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteJSON(message)
		}
	}
}

func (c *Client) readMessage() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[broadcast.client.readMessage] [ERROR] %v", err)
			}
			break
		}
		log.Printf("[broadcast.client.readMessage] recv: %s", message)
	}
}
