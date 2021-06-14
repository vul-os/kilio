package mongo

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lalela-backend/internal/pkg/mongo"
	"lalela-backend/internal/pkg/submissions"
	submissionsStore "lalela-backend/internal/pkg/submissions/store"
)

type store struct {
	collection *mongo.Collection
}

func New(database *mongo.Database) submissionsStore.Store {
	return &store {
		collection: database.Collection("submissions"),
	}
}

func (s *store) CreateOne(request submissionsStore.CreateOneRequest) (*submissionsStore.CreateOneResponse, error) {
	if err := s.collection.CreateOne(request); err != nil {
		log.Error().Err(err).Msg("error creating one job entity")
		return nil, err
	}
	return &submissionsStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request submissionsStore.FindOneRequest) (*submissionsStore.FindOneResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(request.Id)
	var submissionsModel submissions.Submissions

	if err != nil{
		fmt.Println("Invalid id")
		return &submissionsStore.FindOneResponse{}, err
	}
	if err := s.collection.FindOne(&submissionsModel, bson.M{"_id": objectId}); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("error finding one job entity")
			return nil, err
		}
	}

	return &submissionsStore.FindOneResponse{
		FormId: submissionsModel.FormId.Hex(),
		OrgId: submissionsModel.OrganizationId.Hex(),
		Data: submissionsModel.SubmissionData,
	}, nil
}

