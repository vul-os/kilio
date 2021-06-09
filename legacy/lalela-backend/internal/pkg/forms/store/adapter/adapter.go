package adapter

import (
	"github.com/rs/zerolog/log"
	jsonRPCServiceProvider "lalela-backend/internal/pkg/api/jsonRpc/service/provider"
	formsStore "lalela-backend/internal/pkg/forms/store"
	"net/http"
)

type adaptor struct {
	store formsStore.Store
}

func New(
	store formsStore.Store,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		store: store,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return formsStore.FormsServiceProvider
}

func (a *adaptor) CreateOne(r *http.Request, request *formsStore.CreateOneRequest,
	response *formsStore.CreateOneResponse) error {
	result, err := a.store.CreateOne(*request)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	response.Id = result.Id
	return nil
}

func (a *adaptor) FindOne(r *http.Request, request *formsStore.FindOneRequest,
	response *formsStore.FindOneResponse) error {
	model, err := a.store.GetOne(*request)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	response = model
	return nil
}
