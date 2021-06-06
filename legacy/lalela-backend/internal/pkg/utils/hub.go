package utils

import (
	"encoding/json"
	"lalela-backend/internal/pkg/models"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*WsClient]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *WsClient

	// Unregister requests from clients.
	unregister chan *WsClient
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *WsClient),
		unregister: make(chan *WsClient),
		clients:    make(map[*WsClient]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			msg := models.Message{}
			err := json.Unmarshal(message, &msg)
			if err != nil {
				panic(err)
			}
			// only send message to clients that belong to this hub
			for client := range h.clients {
				// only send the message to the people in the room
				if client.room == msg.ChatRoomName {
					// log.Printf("MongoClient room: %s message room: %s \n",
					// 	client.room, msg.ChatRoomName)
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
