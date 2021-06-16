package store

import "lalela-backend/internal/pkg/organizations"

type Store interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	FindOne(FindOneRequest) (*FindOneResponse, error)
}

const OrgsServiceProvider = "Organizations-Store"

const OrgsCreateOneService = OrgsServiceProvider + ".CreateOne"
const OrgsFindOneService = OrgsServiceProvider + ".FindOne"


type CreateOneRequest struct {
	Org organizations.Organizations `json:"organization"`
}

type CreateOneResponse struct {
	Id string `json:"id"`
}

type FindOneRequest struct {
	Identifier string `json:"identifier"`
}

type FindOneResponse struct {
	Org organizations.Organizations `json:"organization"`
}