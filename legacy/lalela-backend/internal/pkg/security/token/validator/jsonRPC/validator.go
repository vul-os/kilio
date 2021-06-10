package jsonRPC

import (
	"github.com/rs/zerolog/log"
	jsonRpcClient "scraperama/internal/pkg/api/jsonRpc/client"
	ybbusJsonRpcClient "scraperama/internal/pkg/api/jsonRpc/client/ybbus"
	"scraperama/internal/pkg/exception"
	tokenValidator "scraperama/internal/pkg/security/token/validator"
	tokenValidatorJSONRPCAdaptor "scraperama/internal/pkg/security/token/validator/adaptor/jsonRPC"
)

type validator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	url, preSharedSecret string,
) tokenValidator.Validator {
	return &validator{
		jsonRpcClient: ybbusJsonRpcClient.New(url, preSharedSecret),
	}
}

func (a *validator) Validate(request tokenValidator.ValidateRequest) (*tokenValidator.ValidateResponse, error) {
	response := new(tokenValidatorJSONRPCAdaptor.ValidateResponse)
	if err := a.jsonRpcClient.JSONRPCRequest(
		tokenValidator.ValidateService,
		nil,
		tokenValidatorJSONRPCAdaptor.ValidateRequest{
			Token: request.Token,
		},
		response,
	); err != nil {
		log.Error().Err(err).Msg("TokenValidator.Validate json rpc")
		return nil, exception.ErrUnexpected{}
	}
	return &tokenValidator.ValidateResponse{Claims: response.Claims.Claims}, nil
}
