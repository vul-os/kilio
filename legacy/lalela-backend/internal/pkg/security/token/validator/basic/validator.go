package token

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gopkg.in/square/go-jose.v2"
	"lalela-backend/internal/pkg/security/claims"
	"lalela-backend/internal/pkg/security/token"
	tokenValidator "lalela-backend/internal/pkg/security/token/validator"
)

type validator struct {
	rsaKeyPair       *rsa.PrivateKey
}

func New(
	rsaKeyPair *rsa.PrivateKey,
) tokenValidator.Validator {
	return &validator{
		rsaKeyPair:       rsaKeyPair,
	}
}

func (v *validator) Validate(request tokenValidator.ValidateRequest) (*tokenValidator.ValidateResponse, error) {
	// Parse the jwt. Successful parse means the received token string is a jwt
	jwtObject, err := jose.ParseSigned(request.Token)
	if err != nil {
		return nil, token.ErrInvalidToken{Reasons: []string{err.Error()}}
	}

	// Verify jwt signature and retrieve json marshalled claims
	// Failure indicates jwt was damaged or tampered with
	jsonClaims, err := jwtObject.Verify(&v.rsaKeyPair.PublicKey)
	if err != nil {
		return nil, token.ErrTokenVerification{Reasons: []string{err.Error()}}
	}

	// unmarshal claims
	var serializedClaims claims.Serialized
	if err := json.Unmarshal(jsonClaims, &serializedClaims); err != nil {
		log.Warn().Err(err).Msg("could not unmarshal claims")
		return nil, err
	}

	// check that claims are not expired
	if serializedClaims.Claims.Expired() {
		return nil, token.ErrInvalidToken{Reasons: []string{"token expired"}}
	}

	// return marshalled claims
	return &tokenValidator.ValidateResponse{Claims: serializedClaims.Claims}, nil
}
