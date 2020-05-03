package models

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Client struct {
	id string
	socket *websocket.Conn
	send chan []byte
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Read(manager *ClientManager){
	defer func() {
		manager.unregister <- c
		_ = c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <-c
			_ = c.socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&Message{
			Sender:    c.id,
			Content:   string(message),
		})
		manager.broadcast <- jsonMessage
	}
}

