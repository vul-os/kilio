package adapter

import (
	"github.com/rs/zerolog/log"
	jsonRPCServiceProvider "lalela-backend/internal/pkg/api/jsonRpc/service/provider"
	userStore "lalela-backend/internal/pkg/users/store"
	"net/http"
	"time"
)


type adaptor struct {
	store userStore.Store
}

func New(
	store userStore.Store,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		store: store,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return userStore.UserServiceProvider
}

type CreateUserRequest struct {
	ID              string 				`json:"id"`
	OrganizationId  string 				`json:"organization_id"`
	FirstName       string             `json:"first_name"`
	LastName        string             `json:"last_name"`
	Password        string             `json:"password"`
	Email           string             `json:"email"`
	RoleID          string             `json:"role_id"`
	ValidationToken string             `json:"validation_token"`
	EmailToken      string             `json:"email_token"`
	RefreshToken    string             `json:"refresh_token"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type FindUserResponse struct {
	FirstName       string             `json:"first_name"`
	LastName        string             `json:"last_name"`
}

func (a *adaptor) CreateUser(r *http.Request, request *CreateUserRequest,
	response *userStore.CreateUserResponse) error {

	result, err := a.store.CreateUser(userStore.CreateUserRequest{
		FirstName: request.FirstName,
		LastName: request.LastName,
		RoleID: request.RoleID,
	})

	if err != nil {
		log.Error().Err(err)
		return err
	}
	response.Id = result.Id
	return nil
}

func (a *adaptor) FindUser(r *http.Request, request *userStore.FindUserRequest,
	response *FindUserResponse) error {

	result, err := a.store.FindUser(*request)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	response = &FindUserResponse{
		FirstName: result.FirstName,
		LastName: result.LastName,

	}
	return nil
}
/*type UserCon struct{}

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

	foundUserId, err := store.CreateUser(args.Email, args.Password)
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
}*/

/*func (t *UserCon) GetOne(r *http.Request, args *GetOneRequest,	reply *GetOneResponse) error {
	canDo, err := auth.Authorize(r, "user", "get")
	if err != nil {
		return err
	}
	if !canDo {
		return nil
	}

	return nil
}*/
