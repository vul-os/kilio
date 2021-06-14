package forms

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)


type Forms struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FormName  string             `json:"form_name"`
	Scheme    interface{}        `json:"scheme"`
	UiScheme  interface{}        `json:"ui_scheme"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}






