package token

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gopkg.in/square/go-jose.v2"
	scraperamaException "scraperama/internal/pkg/exception"
	"scraperama/internal/pkg/security/claims"
	tokenGenerator "scraperama/internal/pkg/security/token/generator"
	"scraperama/pkg/validate/validator"
)

type generator struct {
	tokenSigner jose.Signer
	validator   validator.Validator
}

func New(
	tokenSigner jose.Signer,
	validator validator.Validator,
) tokenGenerator.Generator {
	return &generator{
		tokenSigner: tokenSigner,
		validator:   validator,
	}
}

func (g *generator) GenerateToken(request *tokenGenerator.GenerateTokenRequest) (*tokenGenerator.GenerateTokenResponse, error) {
	if err := g.validator.Validate(request); err != nil {
		return nil, err
	}

	// marshall claims
	claimsPayload, err := json.Marshal(claims.Serialized{
		Claims: request.Claims,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not marshal claims for token")
		return nil, scraperamaException.ErrUnexpected{}
	}

	// sign marshalled payload
	signedObj, err := g.tokenSigner.Sign(claimsPayload)
	if err != nil {
		log.Error().Err(err).Msg("could not sign payload")
		return nil, scraperamaException.ErrUnexpected{}
	}

	// serialize signed object
	signedJWT, err := signedObj.CompactSerialize()
	if err != nil {
		log.Error().Err(err).Msg("could not serialize signed token")
		return nil, scraperamaException.ErrUnexpected{}
	}

	return &tokenGenerator.GenerateTokenResponse{Token: signedJWT}, nil
}
