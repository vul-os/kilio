package basic

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	lalelaException "lalela-backend/internal/pkg/exception"
	"lalela-backend/internal/pkg/mongo"
	lalelaAuthenticator "lalela-backend/internal/pkg/security/authenticator"
	"lalela-backend/internal/pkg/security/casbin"
	"lalela-backend/internal/pkg/security/claims"
	tokenGenerator "lalela-backend/internal/pkg/security/token/generator"
	"lalela-backend/internal/pkg/users"
	usersStore "lalela-backend/internal/pkg/users/store"
	"strings"
	"time"
)

type authenticator struct {
	usersStore        usersStore.Store
	tokenGenerator   tokenGenerator.Generator
	database         *mongo.Database
	casbinEnforcer 	 *casbin.Casbin
}

func New(
	usersStore usersStore.Store,
	tokenGenerator tokenGenerator.Generator,
	database *mongo.Database,
	enforcer 	 *casbin.Casbin,
) lalelaAuthenticator.Authenticator {
	return &authenticator{
		usersStore:       usersStore,
		tokenGenerator:   tokenGenerator,
		database:         database,
		casbinEnforcer:   enforcer,
	}
}

func (a *authenticator) Login(request lalelaAuthenticator.LoginRequest) (*lalelaAuthenticator.LoginResponse, error) {
	var userLoggingIn users.User
	log.Info().Msg(request.Email)
	err := a.database.Collection("users").FindOne(&userLoggingIn, bson.M{"email": request.Email})
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
		fmt.Println(request.Claims, typedClaims.UserID)
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

		serv := strings.Split(request.Service, ".")
		if len(serv) != 2 {
			log.Error().Err(err).Msg("service is not correct")
			return nil, lalelaException.ErrUnauthorized{Reason: "service is not correct"}
		}
		resultEnf, err := a.casbinEnforcer.Enforcer.Enforce(findOneUserResponse.User.ID, request.OrganizationId, serv[0], serv[1])

		if err != nil {
			log.Error().Err(err).Msg("enforcer error")
			return nil, lalelaException.ErrUnauthorized{Reason: "enforcer error: " + err.Error()}
		}
		if resultEnf {
			return &lalelaAuthenticator.AuthenticateServiceResponse{}, nil
		}
	}

	return nil, lalelaException.ErrUnauthorized{Reason: "invalid claims"}
}
