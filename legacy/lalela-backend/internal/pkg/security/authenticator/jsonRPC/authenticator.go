package jsonRPC

import (
	"github.com/rs/zerolog/log"
	jsonRpcClient "lalela-backend/internal/pkg/api/jsonRpc/client"
	ybbusJsonRpcClient "lalela-backend/internal/pkg/api/jsonRpc/client/ybbus"
	lalelaException "lalela-backend/internal/pkg/exception"
	lalelaAuthenticator "lalela-backend/internal/pkg/security/authenticator"
	lalelaAuthenticatorJSONRPCAdaptor "lalela-backend/internal/pkg/security/authenticator/adaptor/jsonRpc"
	"lalela-backend/internal/pkg/security/claims"
)

type authenticator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	url, preSharedSecret string,
) lalelaAuthenticator.Authenticator {
	log.Info().Msg("role json rpc store for: " + url)
	return &authenticator{
		jsonRpcClient: ybbusJsonRpcClient.New(url, preSharedSecret),
	}
}

func (a *authenticator) Login(request lalelaAuthenticator.LoginRequest) (*lalelaAuthenticator.LoginResponse, error) {
	return nil, lalelaException.ErrUnexpected{Reasons: []string{"not implemented"}}
}

func (a *authenticator) AuthenticateService(request lalelaAuthenticator.AuthenticateServiceRequest) (*lalelaAuthenticator.AuthenticateServiceResponse, error) {

	authenticateServiceResponse := new(lalelaAuthenticatorJSONRPCAdaptor.AuthenticateServiceResponse)
	if err := a.jsonRpcClient.JSONRPCRequest(
		lalelaAuthenticator.AuthenticateServiceService,
		nil,
		lalelaAuthenticatorJSONRPCAdaptor.AuthenticateServiceRequest{
			Claims: claims.Serialized{
				Claims: request.Claims,
			},
			Service: request.Service,
		},
		authenticateServiceResponse); err != nil {
		log.Error().Err(err).Msg("auth authenticator jsonrpc authenticateService")
		return nil, err
	}

	return &lalelaAuthenticator.AuthenticateServiceResponse{}, nil
}
