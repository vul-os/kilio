package adapter

import (
	"encoding/json"
	"fmt"
	"lalela-backend/internal/pkg/auth"
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
	return shopStore.ServiceProvider
}

type FormsCon struct{}



func (t *FormsCon) CreateOne(r *http.Request, args *CreateOneRequest,	reply *CreateOneResponse) error {
	canDo, err := auth.Authorize(r, "forms", "create")
	if err != nil {
		return err
	}
	if !canDo {
		return nil
	}

	err = store.CreateOne(args.Name, args.Scheme, args.UiScheme)
	if err != nil {
		return err
	}
	return nil
}

type GetOneRequest struct {
	Id string `json:"id"`
}

type GetOneResponse struct {
	Name string `json:"name"`
	Scheme json.RawMessage `json:"scheme"`
	UiScheme json.RawMessage `json:"ui_scheme"`
}


func (t *FormsCon) GetOne(r *http.Request, args *GetOneRequest, reply *GetOneResponse) error {
	canDo, err := auth.Authorize(r, "forms", "get")
	if err != nil {
		return err
	}
	if !canDo {
		return nil
	}

	model, err := store.GetOne(args.Id)
	if err != nil {
		return err
	}
	fmt.Println(model)
	return nil
}