package middleware

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"lalela-backend/internal/pkg/security/claims"
	tokenValidator "lalela-backend/internal/pkg/security/token/validator"
)

// Authentication Middleware is used to ensure that invoking user is authorised
// i.e. it confirms that they have a valid authentication token
type Authentication struct {
	tokenValidator tokenValidator.Validator
}

func NewAuthentication(
	tokenValidator tokenValidator.Validator,
) *Authentication {

	return &Authentication{
		tokenValidator: tokenValidator,
	}
}

func (a *Authentication) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// an authorisation header with a valid jwt token should be present
		jwt := r.Header.Get("Authorization")
		if jwt == "" {
			log.Error().Msg("no token in authorisation header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// validate token to get serialized claims
		validateResponse, err := a.tokenValidator.Validate(
			tokenValidator.ValidateRequest{
				Token: jwt,
			},
		)
		if err != nil {
			log.Error().Err(err).Msg("token validation failure")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// marshall claims to put into context
		marshalledClaims, err := json.Marshal(claims.Serialized{Claims: validateResponse.Claims})
		if err != nil {
			log.Error().Err(err).Msg("could not marshall claims")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "Claims", marshalledClaims)))
	})
}
