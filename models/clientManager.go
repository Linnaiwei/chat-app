package models

import (
	"encoding/json"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{
				Content:   "/A new socket has been connected",
			})
			manager.send(jsonMessage, conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok{
				close(conn.send)
				delete(manager.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{
					Content:   "/A socket has been disconnected",
				})
				manager.send(jsonMessage, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients{
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

// send message to every client except ignored one
func (manager *ClientManager) send (message []byte, ignore *Client) {
	for client := range manager.clients{
		if client != ignore{
			client.send <- message
		}
	}
}