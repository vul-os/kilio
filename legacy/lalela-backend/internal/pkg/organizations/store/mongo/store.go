package mongo

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	orgsStore "lalela-backend/internal/pkg/organizations/store"
	"lalela-backend/internal/pkg/organizations"
	"lalela-backend/internal/pkg/mongo"
)

type store struct {
	collection *mongo.Collection
}

func New(database *mongo.Database) orgsStore.Store {
	return &store {
		collection: database.Collection("forms"),
	}
}

func (s *store) CreateOne(request orgsStore.CreateOneRequest) (*orgsStore.CreateOneResponse, error) {
	if err := s.collection.CreateOne(request); err != nil {
		log.Error().Err(err).Msg("error creating one job entity")
		return nil, err
	}
	return &orgsStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request orgsStore.FindOneRequest) (*orgsStore.FindOneResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(request.Id)
	var orgsModel organizations.Organizations

	if err != nil{
		fmt.Println("Invalid id")
		return &orgsStore.FindOneResponse{}, err
	}
	if err := s.collection.FindOne(&orgsModel, bson.M{"_id": objectId}); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("error finding one job entity")
			return nil, err
		}
	}

	return &orgsStore.FindOneResponse{
		Name: orgsModel.Name,

	}, nil
}

