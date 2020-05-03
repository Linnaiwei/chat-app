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

func (c *Client) Write() {
	defer func() {
		_ = c.socket.Close()
	}()

	for {
		select {
		case message, ok := <- c.send:
			if !ok {
				_ = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
			}
			_ = c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// Client builder pattern code
type ClientBuilder struct {
	client *Client
}

func NewClientBuilder() *ClientBuilder {
	client := &Client{}
	b := &ClientBuilder{client: client}
	return b
}

func (b *ClientBuilder) Id(id string) *ClientBuilder {
	b.client.id = id
	return b
}

func (b *ClientBuilder) Socket(socket *websocket.Conn) *ClientBuilder {
	b.client.socket = socket
	return b
}

func (b *ClientBuilder) Send(send chan []byte) *ClientBuilder {
	b.client.send = send
	return b
}

func (b *ClientBuilder) Build() (*Client, error) {
	return b.client, nil
}

