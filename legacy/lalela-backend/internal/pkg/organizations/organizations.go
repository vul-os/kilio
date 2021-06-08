package organizations

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)


type Organizations struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  	  string             `json:"organization_name"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

