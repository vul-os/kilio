package mongo

import (
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"lalela-backend/internal/pkg/mongo"
	"lalela-backend/internal/pkg/organizations"
	orgsStore "lalela-backend/internal/pkg/organizations/store"
)

type store struct {
	collection *mongo.Collection
}

func New(database *mongo.Database) orgsStore.Store {
	orgCollection := database.Collection("organizations")
	// setup collection indices
	if err := orgCollection.SetupIndices(
		[]mongoDriver.IndexModel{
			mongo.NewUniqueIndex("name"),
		},
	); err != nil {
		log.Fatal().Err(err).Msg("error setting up user collection indices")
	}
	return &store {
		collection: orgCollection,
	}
}

func (s *store) CreateOne(request orgsStore.CreateOneRequest) (*orgsStore.CreateOneResponse, error) {
	err := s.collection.CreateOne(organizations.Organizations{
		ID:         request.Org.ID,
		Name:       request.Org.Name,
	});
	if  err != nil {
		log.Error().Err(err).Msg("error creating org")
		return nil, err
	}
	return &orgsStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request orgsStore.FindOneRequest) (*orgsStore.FindOneResponse, error) {
	var orgsModel organizations.Organizations
	if err := s.collection.FindOne(&orgsModel, bson.M{"id": request.Identifier}); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("error finding one job entity")
			return nil, err
		}
	}

	return &orgsStore.FindOneResponse{
		Org: orgsModel,
	}, nil
}

