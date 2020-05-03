package models

type ClientManager struct{
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}