package jsonRPC

import (
	"github.com/rs/zerolog/log"
	jsonRPCServiceProvider "lalela-backend/internal/pkg/api/jsonRpc/service/provider"
	scraperamaException "lalela-backend/internal/pkg/exception"
	"lalela-backend/internal/pkg/security/claims"
	"lalela-backend/internal/pkg/users"
	userStore "lalela-backend/internal/pkg/users/store"
	"net/http"
)

type adaptor struct {
	store userStore.Store
}

func New(
	authenticator userStore.Store,
) *adaptor {
	return &adaptor{
		store: authenticator,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return userStore.ServiceProvider
}

type CreateOneRequest struct {
	User users.User `json:"user"`
}

type CreateOneResponse struct {
	User users.User `json:"user"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	if _, err := a.store.CreateOne(
		userStore.CreateOneRequest{
			User: request.User,
		},
	); err != nil {
		return err
	}

	return nil
}

type FindOneRequest struct {
	Identifier string `json:"identifier"`
}

type FindOneResponse struct {
	User users.User `json:"user"`
}

func (a *adaptor) FindOne(r *http.Request, request *FindOneRequest, response *FindOneResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not pass claims for context")
		return scraperamaException.ErrUnexpected{}
	}

	findOneResponse, err := a.store.FindOne(
		userStore.FindOneRequest{
			Claims:     c,
			Identifier: request.Identifier,
		},
	)
	if err != nil {
		return err
	}

	response.User = findOneResponse.User

	return nil
}

type UpdateOneRequest struct {
	User users.User `json:"user"`
}

type UpdateOneResponse struct {
	User users.User `json:"user"`
}

func (a *adaptor) UpdateOne(r *http.Request, request *UpdateOneRequest, response *UpdateOneResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not pass claims for context")
		return scraperamaException.ErrUnexpected{}
	}

	if _, err := a.store.UpdateOne(
		userStore.UpdateOneRequest{
			Claims: c,
			User:   request.User,
		},
	); err != nil {
		return err
	}

	return nil
}
