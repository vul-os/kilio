package generator

import (
	"lalela-backend/internal/pkg/security/claims"
)

type Generator interface {
	GenerateToken(request *GenerateTokenRequest) (*GenerateTokenResponse, error)
}

const ServiceProvider = "Token-Generator"

const GenerateTokenService = ServiceProvider + ".GenerateToken"

type GenerateTokenRequest struct {
	Claims claims.Claims `validate:"required"`
}

type GenerateTokenResponse struct {
	Token string
}
