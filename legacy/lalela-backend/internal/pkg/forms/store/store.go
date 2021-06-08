package store

import (
	"encoding/json"
)

type Store interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	GetOne(GetOneRequest) (*GetOneResponse, error)
}

type CreateOneRequest struct {
	Name string `json:"name"`
	Scheme json.RawMessage `json:"scheme"`
	UiScheme json.RawMessage `json:"ui_scheme"`
}

type CreateOneResponse struct {
	Id string `json:"id"`
}

type GetOneRequest struct {
	Id string `json:"id"`
}

type GetOneResponse struct {
	Name string `json:"name"`
	Scheme interface{} `json:"scheme"`
	UiScheme interface{} `json:"ui_scheme"`
}