package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Chatroom struct {
	Id        primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string               `bson:"name" json:"name"`
	Level     string               `bson:"level" json:"level"`
	Active    string               `bson:"active" json:"active"`
	Timestamp time.Time            `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	Messages  []primitive.ObjectID `bson:"messages,omitempty" json:"messages"`
}

type Message struct {
	MessageId    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Level        int                `bson:"level" json:"level"`
	Text         string             `bson:"text" json:"text"`
	UserName     string             `bson:"name" json:"name"`
	ChatRoomName string             `bson:"room_name" json:"room_name"`
	ChatRoomId   primitive.ObjectID `bson:"chatRoomId,omitempty" json:"chatRoomId,omitempty"`
	Timestamp    time.Time          `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}
