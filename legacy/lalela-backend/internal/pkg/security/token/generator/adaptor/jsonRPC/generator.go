package jsonRPC

import (
	"net/http"
	jsonRPCServiceProvider "lalela-backend/internal/pkg/api/jsonRpc/service/provider"
	"lalela-backend/internal/pkg/security/claims"
	tokenGenerator "lalela-backend/internal/pkg/security/token/generator"
)

type adaptor struct {
	generator tokenGenerator.Generator
}

func New(
	generator tokenGenerator.Generator,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		generator: generator,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return tokenGenerator.GenerateTokenService
}

type GenerateTokenRequest struct {
	Claims claims.Serialized `json:"claims"`
}

type GenerateTokenResponse struct {
	Token string `json:"token"`
}

func (a *adaptor) GenerateToken(r *http.Request, request *GenerateTokenRequest, response *GenerateTokenResponse) error {
	generateResponse, err := a.generator.GenerateToken(
		&tokenGenerator.GenerateTokenRequest{
			Claims: request.Claims.Claims,
		},
	)
	if err != nil {
		return err
	}

	response.Token = generateResponse.Token
	return nil
}
