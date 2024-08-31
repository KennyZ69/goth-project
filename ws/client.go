package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	// Message chan *Message
	Send    chan []byte
	User_id uint
	// Chat_id uint
}

// idk if i actually need this but will find out
type Chat struct {
	Chat_id      uint
	user1_id     uint
	user2_id     uint
	last_message time.Time
}

type Message struct {
	Chat_id    uint      `json:"chat_id"`
	Client_id  uint      `json:"client_id"`
	Username   string    `json:"username"`
	Text       string    `json:"text"`
	created_at time.Time `json:"created_at"`
}

// type WsMessage struct {
// 	Text string      `json:"text"`
// 	Headers interface{} `json:"headers"`
// }

const (
	pongWait   = 60 * time.Second
	maxMsgSize = 512
	pingPeriod = (pongWait * 9) / 10
	writeWait  = 10 * time.Second
)

func serveSocket(w http.ResponseWriter, r *http.Request) {
	currentUser, err := components.GetUserByCookie(r)
	if err != nil {
		log.Println(err)
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{
		Conn:    conn,
		Send:    make(chan []byte, 256),
		User_id: currentUser.Id,
	}

	GlobalHub.Register <- client

	go client.writePump()
	go client.readPump()

}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {

		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			writer, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			writer.Write(msg)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				writer.Write(msg)
			}

			if err := writer.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		GlobalHub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMsgSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, content, err := c.Conn.ReadMessage()
		log.Printf("val: %v\n", string(content))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v\n", err)
			}
			break
		}

		msg := &Message{}
		reader := bytes.NewReader(content)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(msg)
		// err = json.Unmarshal(content, msg)
		if err != nil {
			log.Printf("there was error decoding message content: %v\n", err)
		}

		// want to get the username to be able to display it
		user, err := database.GetUserById(c.User_id)
		if err != nil {
			fmt.Printf("error getting user by id in readPump(): %v\n", err)
		}
		log.Printf("Broadcasting message: %v, from user: %v\n", msg.Text, user.Username)
		GlobalHub.Broadcast <- &Message{
			Text:      msg.Text,
			Client_id: c.User_id,
			Username:  user.Username,
		}
	}
}
