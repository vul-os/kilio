package mongo

import (
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"lalela-backend/internal/pkg/mongo"
	"lalela-backend/internal/pkg/users"
	userStore "lalela-backend/internal/pkg/users/store"
)

type store struct {
	collection *mongo.Collection
}

func New(
	database *mongo.Database,
) userStore.Store {
	// get collection
	userCollection := database.Collection("users")

	// setup collection indices
	if err := userCollection.SetupIndices(
		[]mongoDriver.IndexModel{
			mongo.NewUniqueIndex("id"),
			mongo.NewUniqueIndex("email"),
		},
	); err != nil {
		log.Fatal().Err(err).Msg("error setting up user collection indices")
	}

	return &store{
		collection: userCollection,
	}
}

func (s *store) CreateOne(request userStore.CreateOneRequest) (*userStore.CreateOneResponse, error) {

	err := s.collection.CreateOne(users.User{
		ID:         request.User.ID,
		Name:       request.User.Name,
		Email:      request.User.Email,
		Password:   request.User.Password,
	});
	if  err != nil {
		log.Error().Err(err).Msg("error creating user")
		return nil, err
	}
	return &userStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request userStore.FindOneRequest) (*userStore.FindOneResponse, error) {
	var result users.User
	if err := s.collection.FindOne(&result, bson.M{"id": request.Identifier}); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("finding one user")
			return nil, err
		}
	}
	return &userStore.FindOneResponse{User: result}, nil
}

func (s *store) UpdateOne(request userStore.UpdateOneRequest) (*userStore.UpdateOneResponse, error) {
	if err := s.collection.UpdateOne(request.User, request.User.ID); err != nil {
		log.Error().Err(err).Msg("updating user")
		return nil, err
	}
	return &userStore.UpdateOneResponse{}, nil
}
