package adapter

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	jsonRPCServiceProvider "lalela-backend/internal/pkg/api/jsonRpc/service/provider"
	formsStore "lalela-backend/internal/pkg/forms/store"
	submissionsStore "lalela-backend/internal/pkg/submissions/store"
	"net/http"
)

type adaptor struct {
	mongoStore submissionsStore.Store
}

func New(
	store submissionsStore.Store,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		mongoStore: store,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return submissionsStore.FormsServiceProvider
}

type CreateOneRequest struct {
	FormId string      `json:"form_id"`
	OrgId  string      `json:"org_id"`
	Data   json.RawMessage `json:"data"`
}

type FindOneResponse struct {
	FormId string      `json:"form_id"`
	OrgId  string      `json:"org_id"`
	Data   json.RawMessage `json:"data"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest,
	response *formsStore.CreateOneResponse) error {

	var data interface{}
	if err := json.Unmarshal(request.Data, &data); err != nil {
		return err
	}

	result, err := a.mongoStore.CreateOne(submissionsStore.CreateOneRequest{
		FormId: request.FormId,
		OrgId: request.OrgId,
		Data: request.Data,
	})

	if err != nil {
		log.Error().Err(err)
		return err
	}
	response.Id = result.Id
	return nil
}

func (a *adaptor) FindOne(r *http.Request, request *submissionsStore.FindOneRequest,
	response *FindOneResponse) error {

	result, err := a.mongoStore.FindOne(*request)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	data, err := json.Marshal(result.Data)
	if err != nil {
		log.Error().Err(err).Msg("find one json marshalling error")
		return err
	}
	response = &FindOneResponse{
		FormId: result.FormId,
		OrgId: result.OrgId,
		Data: data,
	}
	return nil
}
