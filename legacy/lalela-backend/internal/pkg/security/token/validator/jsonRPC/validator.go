package jsonRPC

import (
	"github.com/rs/zerolog/log"
	jsonRpcClient "lalela-backend/internal/pkg/api/jsonRpc/client"
	ybbusJsonRpcClient "lalela-backend/internal/pkg/api/jsonRpc/client/ybbus"
	"lalela-backend/internal/pkg/exception"
	tokenValidator "lalela-backend/internal/pkg/security/token/validator"
	tokenValidatorJSONRPCAdaptor "lalela-backend/internal/pkg/security/token/validator/adaptor/jsonRPC"
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
	return &tokenValidator.ValidateResponse{Claims: response.Claims}, nil
}
