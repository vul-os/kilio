package middleware

import (
	"fmt"
	"github.com/rs/zerolog/log"
	lalelaAuthenticator "lalela-backend/internal/pkg/security/authenticator"
	"lalela-backend/internal/pkg/security/claims"
	"net/http"
	//"lalela-backend/internal/pkg/security/claims"
)

// Authorisation middleware is used to ensure that the invoking user has
// permission to access the service being called
type Authorisation struct {
	authenticator lalelaAuthenticator.Authenticator
}

func NewAuthorisation(
	authenticator lalelaAuthenticator.Authenticator,
) *Authorisation {

	return &Authorisation{
		authenticator: authenticator,
	}
}

func (a *Authorisation) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get json rpc method
		method, org, err := getMethodAndOrg(r)
		if err != nil {
			log.Error().Err(err).Msg("cannot get jsonrpc method")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//// parse claims from request context
		userClaims, err := claims.ParseClaimsFromContext(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("unable to get claims from context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println(method, org, userClaims)
		//// validate service access
		if _, err := a.authenticator.AuthenticateService(
			lalelaAuthenticator.AuthenticateServiceRequest{
				Claims:  userClaims,
				Service: method,
				OrganizationId: org,
			},
		); err != nil {
			log.Error().Err(err).Msg("unauthorized to access service")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
