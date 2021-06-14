package mongo

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lalela-backend/internal/pkg/mongo"
	"lalela-backend/internal/pkg/users"
	userStore "lalela-backend/internal/pkg/users/store"
	"time"
)

type store struct {
	collection *mongo.Collection
}


func New(database *mongo.Database) userStore.Store {
	return &store {
		collection: database.Collection("users"),
	}
}

func (s *store) CreateUser(request userStore.CreateUserRequest) (*userStore.CreateUserResponse, error) {
	if err := s.collection.CreateOne(request); err != nil {
		log.Error().Err(err).Msg("error creating one job entity")
		return nil, err
	}
	return &userStore.CreateUserResponse{}, nil
}

func (s *store) FindUser(request userStore.FindUserRequest) (*userStore.FindUserResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(request.Id)
	var userModel users.User

	if err != nil{
		fmt.Println("Invalid id", time.Now())
		return &userStore.FindUserResponse{}, err
	}
	if err := s.collection.FindOne(&userModel, bson.M{"_id": objectId}); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("error finding user")
			return nil, err
		}
	}

	return &userStore.FindUserResponse{
		FirstName: userModel.FirstName,
		LastName: userModel.LastName,
		//RoleID: userModel.RoleID,
	}, nil
}
