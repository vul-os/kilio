package store

import (
	"lalela-backend/internal/pkg/security/claims"
	"lalela-backend/internal/pkg/users"
)

type Store interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	FindOne(FindOneRequest) (*FindOneResponse, error)
	UpdateOne(UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "User-Store"

const git  = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"
const FindManyService = ServiceProvider + ".FindMany"
const UpdateOneService = ServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	User users.User
}

type CreateOneResponse struct {
}

type FindOneRequest struct {
	Claims     claims.Claims
	Identifier string
}

type FindOneResponse struct {
	User users.User
}

type UpdateOneRequest struct {
	Claims claims.Claims
	User   users.User
}

type UpdateOneResponse struct {
}
