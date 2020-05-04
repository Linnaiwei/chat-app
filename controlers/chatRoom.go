package controlers

import (
	"chat-app/models"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

var localManager interface{}

func ChatRoomHandler(manager *models.ClientManager) {
	localManager = manager
	http.HandleFunc("/chat-room/", ChatRoomPage)
}

func ChatRoomPage(res http.ResponseWriter, req *http.Request)  {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}).Upgrade(res, req, nil)
	if err != nil{
		http.NotFound(res, req)
		return
	}
	client, _ := models.NewClientBuilder().
		Id(uuid.NewV4().String()).
		Socket(conn).
		Send(make(chan []byte)).
		Build()
	manager, ok := localManager.(*models.ClientManager)
	if ok {
		manager.SetRegister(client)
	} else {
		log.Fatal("type assertion is failed")
	}
	go client.Read(manager)
	go client.Write()
}
