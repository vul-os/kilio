package adapter

import (
	"encoding/json"
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


type CreateOneRequest struct {
	Name string `json:"name"`
	Scheme json.RawMessage `json:"scheme"`
	UiScheme json.RawMessage `json:"ui_scheme"`
}

type FindOneResponse struct {
	Name string `json:"name"`
	Scheme json.RawMessage `json:"scheme"`
	UiScheme json.RawMessage `json:"ui_scheme"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest,
	response *formsStore.CreateOneResponse) error {

	// TODO: fix this ugly bs
	var scheme interface{}
	var uiScheme interface{}
	if err := json.Unmarshal(request.Scheme, &scheme); err != nil {
		return err
	}
	if err := json.Unmarshal(request.UiScheme, &uiScheme); err != nil {
		return err
	}

	result, err := a.store.CreateOne(formsStore.CreateOneRequest{
		Name: request.Name,
		Scheme: scheme,
		UiScheme: uiScheme,
	})

	if err != nil {
		log.Error().Err(err)
		return err
	}
	response.Id = result.Id
	return nil
}

func (a *adaptor) FindOne(r *http.Request, request *formsStore.FindOneRequest,
	response *FindOneResponse) error {

	result, err := a.store.FindOne(*request)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	scheme, _ := json.Marshal(result.Scheme)
	uischeme, _ := json.Marshal(result.UiScheme)
	response = &FindOneResponse{
		Name: result.Name,
		Scheme: scheme,
		UiScheme: uischeme,
	}
	return nil
}
