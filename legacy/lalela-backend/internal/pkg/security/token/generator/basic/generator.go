package token

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gopkg.in/square/go-jose.v2"
	scraperamaException "lalela-backend/internal/pkg/exception"
	"lalela-backend/internal/pkg/security/claims"
	tokenGenerator "lalela-backend/internal/pkg/security/token/generator"
)

type generator struct {
	tokenSigner jose.Signer
}

func New(
	tokenSigner jose.Signer,
) tokenGenerator.Generator {
	return &generator{
		tokenSigner: tokenSigner,
	}
}

func (g *generator) GenerateToken(request *tokenGenerator.GenerateTokenRequest) (*tokenGenerator.GenerateTokenResponse, error) {

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
