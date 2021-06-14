package jsonRPC

import (
	"github.com/rs/zerolog/log"
	jsonRpcClient "lalela-backend/internal/pkg/api/jsonRpc/client"
	ybbusJsonRpcClient "lalela-backend/internal/pkg/api/jsonRpc/client/ybbus"
	tokenGenerator "lalela-backend/internal/pkg/security/token/generator"
	tokenGeneratorJSONRPCAdaptor "lalela-backend/internal/pkg/security/token/generator/adaptor/jsonRPC"
)

type generator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	url, preSharedSecret string,
) tokenGenerator.Generator {
	return &generator{
		jsonRpcClient: ybbusJsonRpcClient.New(url, preSharedSecret),
	}
}

func (a *generator) GenerateToken(request *tokenGenerator.GenerateTokenRequest) (*tokenGenerator.GenerateTokenResponse, error) {
	generateResponse := new(tokenGeneratorJSONRPCAdaptor.GenerateTokenResponse)
	if err := a.jsonRpcClient.JSONRPCRequest(
		tokenGenerator.GenerateTokenService,
		request.Claims,
		request,
		generateResponse,
	); err != nil {
		log.Error().Err(err).Msg("token jsonrpc generator generate")
		return nil, err
	}
	return &tokenGenerator.GenerateTokenResponse{Token: generateResponse.Token}, nil
}
