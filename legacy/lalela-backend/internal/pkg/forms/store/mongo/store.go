package mongo

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"lalela-backend/internal/pkg/forms"
	formsStore "lalela-backend/internal/pkg/forms/store"
	"lalela-backend/internal/pkg/mongo"
)

type store struct {
	collection *mongo.Collection
}

func New(database *mongo.Database) formsStore.Store {
	formCollection := database.Collection("forms")
	// setup collection indices
	if err := formCollection.SetupIndices(
		[]mongoDriver.IndexModel{
			mongo.NewUniqueIndex("FormName"),
		},
	); err != nil {
		log.Fatal().Err(err).Msg("error setting up form collection indices")
	}
	return &store {
		collection: formCollection,
	}
}

func (s *store) CreateOne(request formsStore.CreateOneRequest) (*formsStore.CreateOneResponse, error) {
	if err := s.collection.CreateOne(request); err != nil {
		log.Error().Err(err).Msg("error creating one job entity")
		return nil, err
	}
	return &formsStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request formsStore.FindOneRequest) (*formsStore.FindOneResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(request.Id)
	var formsModel forms.Forms

	if err != nil{
		fmt.Println("Invalid id")
		return &formsStore.FindOneResponse{}, err
	}
	if err := s.collection.FindOne(&formsModel, bson.M{"_id": objectId}); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("error finding one job entity")
			return nil, err
		}
	}

	return &formsStore.FindOneResponse{
		Name: formsModel.FormName,
		Scheme: formsModel.Scheme,
		UiScheme: formsModel.UiScheme,
	}, nil
}

