package models

import "github.com/dgrijalva/jwt-go"

//// AuthCon
type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginResponse struct {
	Jwt  string `json:"jwt,omitempty"`
	User User   `json:"user,omitempty"`
}

type AuthResetRequest struct {
	Email string `json:"email"`
}

type AuthUpdatePasswordRequest struct {
	ResetToken string `json:"resetToken"`
	Password   string `json:"password"`
}

type AuthUpdatePasswordAdminRequest struct {
	Email string `json:"email"`
	User  User
}

type AuthUpdatePasswordAdminResponse struct {
	Messages []string `json:"message"`
}

type AuthLoginJWTRequest struct {
	Jwt string `json:"email"`
}

type AuthLoginJWTResponse struct {
	User User `json:"user,omitempty"`
}

//Token struct declaration
type Token struct {
	UserID uint   `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	RoleID string `json:"roleId"`
	Type   string `json:"type"`
	*jwt.StandardClaims
}