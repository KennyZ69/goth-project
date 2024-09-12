package ws

import (
	"gothstarter/layouts/components"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SocketHandler(w http.ResponseWriter, r *http.Request) error {
	currentUser, err := components.GetUserByCookie(r)
	if err != nil {
		log.Println(err)
		return err
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	client := &Client{
		Conn:    conn,
		Send:    make(chan []byte, 256),
		User_id: currentUser.Id,
	}

	GlobalHub.Register <- client

	go client.writePump()
	go client.readPump(r)

	return nil
}
