package adapter

import (
	"lalela-backend/internal/pkg/auth"
	"lalela-backend/internal/pkg/users"
	"lalela-backend/internal/pkg/users/store"
	"net/http"
)

type UserCon struct{}

type CreateOneRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	OrganizationId string `json:"organization_id"`
}

type CreateOneResponse struct {
	Response string `json:"response"`
	Id string `json:"id"`
}

type GetOneRequest struct {
}

type GetOneResponse struct {
	User users.User
}

type GetManyRequest struct {
}

type GetManyResponse struct {
	User []users.User
}


func (t *UserCon) CreateOne(r *http.Request, args *CreateOneRequest,	reply *CreateOneResponse) error {
	if args.OrganizationId != "" {
		canDo, err := auth.Authorize(r, "user", "create")
		if err != nil {
			return err
		}
		if !canDo {
			return nil
		}
	}

	foundUserId, err := store.CreateOne(args.Email, args.Password)
	if err != nil {
		return err
	}
	reply.Response = "ok"
	reply.Id = foundUserId
	return nil
}



func (t *UserCon) GetMany(r *http.Request, args *GetManyRequest,	reply *GetManyResponse) error {
	canDo, err := auth.Authorize(r, "user", "get")
	if err != nil {
		return err
	}
	if !canDo {
		return nil
	}
	return nil
}

type GetAllUsersRequest struct {
	OrganizationId string `json:"organization_id"`
}

func (t *UserCon) GetOne(r *http.Request, args *GetOneRequest,	reply *GetOneResponse) error {
	canDo, err := auth.Authorize(r, "user", "get")
	if err != nil {
		return err
	}
	if !canDo {
		return nil
	}

	return nil
}
