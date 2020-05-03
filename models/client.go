package models

import "github.com/gorilla/websocket"

type Client struct {
	id string
	socket *websocket.Conn
	send chan []byte
}
