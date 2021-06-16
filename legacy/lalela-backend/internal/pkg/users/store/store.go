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

const UserServiceProvider = "User-Store"

const CreateOneService = UserServiceProvider + ".CreateOne"
const FindOneService = UserServiceProvider + ".FindOne"
const FindManyService = UserServiceProvider + ".FindMany"
const UpdateOneService = UserServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	User users.User
}

type CreateOneResponse struct {
	ID string
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
