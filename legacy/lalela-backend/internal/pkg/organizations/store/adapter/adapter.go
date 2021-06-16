package adapter

import (
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	jsonRPCServiceProvider "lalela-backend/internal/pkg/api/jsonRpc/service/provider"
	"lalela-backend/internal/pkg/organizations"
	orgsStore "lalela-backend/internal/pkg/organizations/store"
	"net/http"
)


type adaptor struct {
	mongoStore orgsStore.Store
}

func New(
	store orgsStore.Store,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		mongoStore: store,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return orgsStore.OrgsServiceProvider
}


type CreateOneRequest struct {
	Name string `json:"name"`
}

type FindOneResponse struct {
	Org organizations.Organizations `json:"organization"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest,
	response *orgsStore.CreateOneResponse) error {

	result, err := a.mongoStore.CreateOne(orgsStore.CreateOneRequest{
		Org: organizations.Organizations{
			ID:   uuid.NewV4().String(),
			Name: request.Name,
		},
	})

	if err != nil {
		log.Error().Err(err)
		return err
	}
	response.Id = result.Id
	return nil
}

func (a *adaptor) FindOne(r *http.Request, request *orgsStore.FindOneRequest,
	response *FindOneResponse) error {
	findOneResponse, err := a.mongoStore.FindOne(
		orgsStore.FindOneRequest{
			Identifier: request.Identifier,
		},
	)
	if err != nil {
		return err
	}

	response.Org = findOneResponse.Org

	return nil
}
//func (a *adaptor) FindOne(r *http.Request, request *orgsStore.FindOneRequest,
//	response *FindOneResponse) error {
//
//	var org organizations.Organizations
//	result, err := a.mongoStore.FindOne(request)
//	if err != nil {
//		log.Error().Err(err)
//		return err
//	}
//
//	response = &FindOneResponse{
//		,
//	}
//	return nil
//}
