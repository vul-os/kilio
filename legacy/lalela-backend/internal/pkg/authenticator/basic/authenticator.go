package basic

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	lalelaAuthenticator "lalela-backend/internal/pkg/authenticator"
	lalelaException "lalela-backend/internal/pkg/exception"
	"lalela-backend/internal/pkg/mongo"
	"lalela-backend/internal/pkg/security/claims"
	tokenGenerator "lalela-backend/internal/pkg/security/token/generator"
	"lalela-backend/internal/pkg/users"
	usersStore "lalela-backend/internal/pkg/users/store"
	"time"
)

type authenticator struct {
	usersStore        usersStore.Store
	tokenGenerator   tokenGenerator.Generator
	database         *mongo.Database
}

func New(
	usersStore usersStore.Store,
	tokenGenerator tokenGenerator.Generator,
	database *mongo.Database,
) lalelaAuthenticator.Authenticator {
	return &authenticator{
		usersStore:       usersStore,
		tokenGenerator:   tokenGenerator,
		database:         database,
	}
}

func (a *authenticator) Login(request lalelaAuthenticator.LoginRequest) (*lalelaAuthenticator.LoginResponse, error) {
	var userLoggingIn users.User
	err := a.database.Collection("user").FindOne(&userLoggingIn, request.Email)
	if err != nil {
		log.Error().Err(err).Msg("retrieving user for log in")
		return nil, err
	}

	// check password is correct
	if err := bcrypt.CompareHashAndPassword(userLoggingIn.Password, []byte(request.Password)); err != nil {
		log.Error().Err(err).Msg("invalid password login")
		return nil, err
	}

	// generate login claims
	generateTokenResponse, err := a.tokenGenerator.GenerateToken(
		&tokenGenerator.GenerateTokenRequest{
			Claims: claims.Login{
				UserID:         userLoggingIn.ID,
				ExpirationTime: time.Now().Add(time.Hour * 1).UTC().Unix(),
			},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("generating token")
		return nil, lalelaException.ErrUnexpected{}
	}

	return &lalelaAuthenticator.LoginResponse{
		JWT: generateTokenResponse.Token,
	}, nil
}

func (a *authenticator) AuthenticateService(request lalelaAuthenticator.AuthenticateServiceRequest) (*lalelaAuthenticator.AuthenticateServiceResponse, error) {
	switch typedClaims := request.Claims.(type) {
	case claims.Login:
		// try and retrieve user that owns claims
		findOneUserResponse, err := a.usersStore.FindOne(
			usersStore.FindOneRequest{
				Claims:     request.Claims,
				Identifier: typedClaims.UserID,
			},
		)
		if err != nil {
			log.Error().Err(err).Msg("could not retrieve user")
			return nil, lalelaException.ErrUnauthorized{Reason: "could not retrieve user: " + err.Error()}
		}
		fmt.Println(findOneUserResponse)
		// todo: do auth
		// create criterion to retrieve user's roles
		//roleCriteria := text.List{
		//	Field: "id",
		//	List:  make([]string, 0),
		//}
		//for _, roleID := range findOneUserResponse.User.RoleIDs {
		//	roleCriteria.List = append(roleCriteria.List, roleID.String())
		//}
		//roleFindManyResponse, err := a.roleStore.FindMany(
		//	roleStore.FindManyRequest{
		//		Criteria: []criterion.Criterion{
		//			roleCriteria,
		//			text.Exact{
		//				Field: "permissions",
		//				Text:  request.Service,
		//			},
		//		},
		//	},
		//)
		if err != nil {
			return nil, lalelaException.ErrUnauthorized{Reason: "could not retrieve roles: " + err.Error()}
		}

		// if any roles match this criteria then the user has access to this service
		//if roleFindManyResponse.Total > 0 {
		//	return &lalelaAuthenticator.AuthenticateServiceResponse{}, nil
		//}
		return nil, lalelaException.ErrUnauthorized{Reason: "no permission"}
	}

	return nil, lalelaException.ErrUnauthorized{Reason: "invalid claims"}
}
