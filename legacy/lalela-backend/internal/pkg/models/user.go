package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName       *string            `json:"first_name"`
	LastName        *string            `json:"last_name"`
	Password        *string            `json:"password"`
	Email           *string            `json:"email"`
	UserType        *string            `json:"user_type"`
	ValidationToken *string            `json:"validation_token"`
	EmailToken      *string            `json:"email_token"`
	RefreshToken    *string            `json:"refresh_token"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}
