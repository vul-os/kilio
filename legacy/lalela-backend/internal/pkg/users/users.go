package users

import (
	"time"
)

type User struct {
	FirstName       string             `json:"first_name"`
	LastName        string             `json:"last_name"`
	Password        string             `json:"password"`
	Email           string             `json:"email"`
	RoleID          string             `json:"role_id"`
	ValidationToken string             `json:"validation_token"`
	EmailToken      string             `json:"email_token"`
	RefreshToken    string             `json:"refresh_token"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}


