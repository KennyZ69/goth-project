package ws

import (
	"bytes"
	"html/template"
	"log"
	"sync"
)

type Hub struct {
	Connections map[*Client]bool

	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client

	Messages []*Message

	sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Connections: make(map[*Client]bool),
		Broadcast:   make(chan *Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
	}
}

var GlobalHub *Hub = NewHub()

func (h *Hub) Run() {

	for {
		select {
		case client := <-h.Register:
			h.Lock()
			h.Connections[client] = true
			h.Unlock()

			log.Printf("client was registered: %v", client.User_id)

			for _, msg := range h.Messages {
				client.Send <- getMsgByteTempl(msg)
			}

		case client := <-h.Unregister:
			h.Lock()
			if _, ok := h.Connections[client]; ok {
				close(client.Send)
				log.Printf("client unregistered %v", client.User_id)
				delete(h.Connections, client)
			}
			h.Unlock()
		case msg := <-h.Broadcast:
			h.RLock()
			h.Messages = append(h.Messages, msg)

			for client := range h.Connections {
				select {
				case client.Send <- getMsgByteTempl(msg):
					log.Printf("Message sent from client: %v, content: %s\n", client.User_id, msg.Text)
				default:
					close(client.Send)
					delete(h.Connections, client)
				}
			}
			h.RUnlock()
		}
	}
}

func getMsgByteTempl(message *Message) []byte {
	tmpl, err := template.ParseFiles("html/message.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}

	// Render the template with the message as data.
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, message)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}

	return renderedMessage.Bytes()
}
