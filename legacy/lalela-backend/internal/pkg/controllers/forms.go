package controllers

import (
	"encoding/json"
	"fmt"
	"lalela-backend/internal/pkg/services"
	"net/http"
)

type FormsCon struct{}

type FormCreateRequest struct {
	Name string `json:"name"`
	Scheme json.RawMessage `json:"scheme"`
	UiScheme json.RawMessage `json:"ui_scheme"`
}

type FormCreateResponse struct {
	Id string `json:"id"`
}


func (t *FormsCon) CreateForm(r *http.Request, args *FormCreateRequest,	reply *FormCreateResponse) error {
	canDo, err := services.Authorize(r, "global", "forms", "create")
	if err != nil {
		return err
	}
	if !canDo {
		return nil
	}

	err = services.CreateForm(args.Name, args.Scheme, args.UiScheme)
	if err != nil {
		return err
	}
	return nil
}


type FormGetRequest struct {
	Id string `json:"id"`
}


func (t *FormsCon) GetForm(r *http.Request, args *FormGetRequest, reply *FormCreateRequest) error {
	canDo, err := services.Authorize(r, "global", "forms", "get")
	if err != nil {
		return err
	}
	if !canDo {
		return nil
	}

	model, err := services.GetForm(args.Id)
	if err != nil {
		return err
	}
	fmt.Println(model)
	return nil
}




