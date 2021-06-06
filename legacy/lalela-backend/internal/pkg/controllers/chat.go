package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/utils"
	"log"
	"net/http"
	"time"
)




// mem
var users map[string]models.User

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

/**
 * Load the previous messages from this channel from the database
 */
func channelHistory(w http.ResponseWriter, r *http.Request) {
	var room Chatroom
	// find the chatroom at this request
	err := Mongo.Chatrooms.Find(bson.M{"name": vars["channel"]}).One(&room)
	if err != nil { // channel not found
		log.Printf("Creating new channel: %s ...", vars["channel"])
		// create new channel
		room.Id = bson.NewObjectId()
		room.Name = vars["channel"]
		room.Level = "0"
		room.Active = "true"
		err := Mongo.Chatrooms.Insert(room)
		if err != nil {
			log.Println(err)
		} else {
			// new welcome message for the room
			welcomeMessage := Message{
				MessageId:    bson.NewObjectId(),
				Text:         "Welcome to the new " + vars["channel"] + " chat",
				ChatRoomName: vars["channel"],
				UserName:     "Moderator",
				ChatRoomId:   room.Id,
				Timestamp:    time.Now(),
				Level:        1, // level = power
			}
			room.Messages = append(room.Messages, welcomeMessage.MessageId)
			// insert the new welcome message into the messages
			// collection, with this chatroom id and the user id
			err = Mongo.Messages.Insert(welcomeMessage)
			if err != nil {
				panic(err) // error inserting
			}
		}
	}
	// initialize a slice of size messageAmount to store the messages
	var messageSlice []Message
	// find the last 150 messages in the room
	err = Mongo.Messages.Find(
		bson.M{"chatRoomId": room.Id}).Sort(
		"-timestamp").Limit(150).All(&messageSlice)
	js, err := json.Marshal(messageSlice)
	if err != nil {
		panic(err)
	}
	// serve
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/**
 * Channel to save messages to the database
 */
func saveMessages(m *chan []byte) {
	for {
		message, ok := <-*m
		if !ok {
			log.Println("Error when trying to save")
			return
		}
		saveMessage(&message)
	}
}

func saveMessage(msg *[]byte) {
	message := models.Message{}
	err := json.Unmarshal(*msg, &message)
	message.MessageId = bson.NewObjectId()
	message.Timestamp = time.Now()
	var room Chatroom
	// find the chatroom at this request
	err = Mongo.Chatrooms.Find(bson.M{"name": message.ChatRoomName}).One(&room)
	if err != nil { // channel not found
		// create new channel
		room.Name = message.ChatRoomName
		room.Level = "0"
		room.Active = "true"
		room.Id = bson.NewObjectId()
		err := Mongo.Chatrooms.Insert(room)
		if err != nil {
			log.Println(err)
		} else {
			room.Messages = append(room.Messages, message.MessageId)
		}
	}
	// construct the new message
	message.ChatRoomId = room.Id
	// insert the message into the messages collection, with this chatroom
	// and the user id
	err = Mongo.Messages.Insert(message)
	if err != nil {
		log.Println(err)
		// panic(err) // error inserting
	}
	var messageSlice []models.Message
	var bsonMessageSlice []bson.ObjectId
	// find all the messages that have this room as chatRoomId
	err = Mongo.Messages.Find(
		bson.M{"chatRoomId": room.Id}).Sort("-timestamp").All(&messageSlice)
	if err != nil {
		panic(err)
	}
	if len(messageSlice) > 0 {
		if err != nil {
			log.Println(err)
		}
		// if there is no messages it won't enter the loop
		for i := 0; i < len(messageSlice); i++ {
			bsonMessageSlice = append(bsonMessageSlice, messageSlice[i].MessageId)
		}
	}
	// append the new message
	bsonMessageSlice = append(bsonMessageSlice, message.MessageId)
	// update the room with the new messsage
	err = Mongo.Chatrooms.Update(bson.M{"_id": room.Id},
		bson.M{"$set": bson.M{"messages": bsonMessageSlice}})
	if err != nil {
		panic(err)
	}
}


// serveWs handles websocket requests from the peer.
func serveWs(hub *utils.Hub, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &utils.WsClient{
		Room: vars["channel"],
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
		Save: &MessageChannel,
	}
	client.Hub.register <- client
	// one goroutine for each client for reading and another for sending
	// messages to and from the hub to the WebSocket
	go client.writePump()
	go client.readPump()
}




