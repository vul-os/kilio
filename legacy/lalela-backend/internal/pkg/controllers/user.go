package controllers

import (
	"lalela-backend/internal/pkg/services"
	"net/http"
)

type UserCon struct{}

type UserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Response string `json:"response"`
	Id string `json:"id"`
}

func (t *UserCon) LoginJWT(r *http.Request, args *UserRequest,	reply *UserResponse) error {

	return nil
}


// https://github.com/Joojo7/user-athentication-golang/blob/master/controllers/userController.go
func (t *UserCon) LoginCredentials(r *http.Request, args *UserRequest,	reply *UserResponse) error {
	foundUserId, err := services.LoginUserCredentials(args.Email, args.Password)
	if err != nil {
		return err
	}
	reply.Response = "ok"
	reply.Id = foundUserId
	return nil
}

type RegisterUserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	OrganizationId string `json:"organization_id"`
}

func (t *UserCon) RegisterUser(r *http.Request, args *RegisterUserRequest,	reply *UserResponse) error {
	if args.OrganizationId != "" {
		canDo, err := services.Authorize(r, args.OrganizationId, "user", "create")
		if err != nil {
			return err
		}
		if !canDo {
			return nil
		}
	}

	foundUserId, err := services.CreateUser(args.Email, args.Password)
	if err != nil {
		return err
	}
	reply.Response = "ok"
	reply.Id = foundUserId
	return nil
}

type GetUserRequest struct {
	id string `json:"id"`
	OrganizationId string `json:"organization_id"`
}

func (t *UserCon) GetUser(r *http.Request, args *GetUserRequest,	reply *UserResponse) error {
	canDo, err := services.Authorize(r, args.OrganizationId, "user", "get")
	if err != nil {
		return err
	}
	if !canDo {
		return nil
	}

	return nil
}

type GetAllUsersInOrgRequest struct {
	OrganizationId string `json:"organization_id"`
}

func (t *UserCon) GetAllUsersInOrg(r *http.Request, args *GetAllUsersInOrgRequest,	reply *UserResponse) error {
	canDo, err := services.Authorize(r, args.OrganizationId, "users", "get")
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

func (t *UserCon) GetAllUsers(r *http.Request, args *GetAllUsersInOrgRequest,	reply *UserResponse) error {
	canDo, err := services.Authorize(r, "global", "users", "create")
	if err != nil {
		return err
	}
	if !canDo {
		return nil
	}
	return nil
}