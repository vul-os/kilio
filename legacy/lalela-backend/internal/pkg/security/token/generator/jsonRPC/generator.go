package jsonRPC

import (
	"github.com/rs/zerolog/log"
	jsonRpcClient "scraperama/internal/pkg/api/jsonRpc/client"
	ybbusJsonRpcClient "scraperama/internal/pkg/api/jsonRpc/client/ybbus"
	"scraperama/internal/pkg/security/claims"
	tokenGenerator "scraperama/internal/pkg/security/token/generator"
	tokenGeneratorJSONRPCAdaptor "scraperama/internal/pkg/security/token/generator/adaptor/jsonRPC"
	"scraperama/pkg/validate/validator"
)

type generator struct {
	jsonRpcClient jsonRpcClient.Client
	validator     validator.Validator
}

func New(
	url, preSharedSecret string,
	validator validator.Validator,
) tokenGenerator.Generator {
	return &generator{
		jsonRpcClient: ybbusJsonRpcClient.New(url, preSharedSecret),
		validator:     validator,
	}
}

func (a *generator) GenerateToken(request *tokenGenerator.GenerateTokenRequest) (*tokenGenerator.GenerateTokenResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	generateResponse := new(tokenGeneratorJSONRPCAdaptor.GenerateTokenResponse)
	if err := a.jsonRpcClient.JSONRPCRequest(
		tokenGenerator.GenerateTokenService,
		tokenGeneratorJSONRPCAdaptor.GenerateTokenRequest{
			Claims: claims.Serialized{
				Claims: request.Claims,
			},
		},
		generateResponse,
	); err != nil {
		log.Error().Err(err).Msg("token jsonrpc generator generate")
		return nil, err
	}
	return &tokenGenerator.GenerateTokenResponse{Token: generateResponse.Token}, nil
}
