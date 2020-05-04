package main

import (
	"chat-app/controlers"
	"chat-app/models"
	"net/http"
)


func main () {
	//spin up a global ClientManager
	manager := models.NewClientManager()
	go manager.Start()
	controlers.ChatRoomHandler(manager)
	_ = http.ListenAndServe(":1234", nil)
}
