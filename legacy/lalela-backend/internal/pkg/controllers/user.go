package controllers

import "cog-analytics-engine-go/internal/pkg/models"

//// UserCon
type UsersGetRequest struct {
	Email string `json:"email"`
}

type UsersGetResponse struct {
	Users []models.User
}

type UserGetRequest struct {
	Email string `json:"email"`
	Id    uint   `json:"id"`
}

type UserGetResponse struct {
	User models.User
}

type UserIsAdminRequest struct {
	Email string `json:"email"`
}

type UserIsAdminResponse struct {
	Admin bool
}

type UserGetResponse2 struct {
	User       models.UserGetReport
	Dashboards map[string][]string
}

type UserAddRequest struct {
	Email string `json:"email"`
	Users []models.User
}

type UserAddResponse struct {
	Messages []string `json:"message"`
}

type UserUpdateRequest struct {
	Email string `json:"email"`
	Users []models.User
}

type UserUpdateRequestSingle struct {
	Email string `json:"email"`
	User  models.User
}

type UserUpdateResponse struct {
	Messages []string `json:"message"`
}

type UserDeleteRequest struct {
	Email string `json:"email"`
	User  models.User
}

type UserDeleteResponse struct {
	Messages []string `json:"message"`
}

type UserRoleUpdateRequest struct {
	Email   string `json:"email"`
	User    models.User
	Roles   []string
	GroupId int
}

type UserRoleUpdateResponse struct {
	Messages []string `json:"message"`
}
