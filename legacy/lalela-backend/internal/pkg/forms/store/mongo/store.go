package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"lalela-backend/internal/pkg/forms"
	formsStore "lalela-backend/internal/pkg/forms/store"
	"time"
)

type store struct {
	collection *mongo.Collection
}

func New(database *mongo.Database) formsStore.Store {
	return &store {
		collection: database.Collection("forms"),
	}
}
func (s *store) CreateOne(request formsStore.CreateOneRequest) (*formsStore.CreateOneResponse, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
	var model forms.Forms

	var sch interface{}
	var uiSch interface{}

	if err := json.Unmarshal(request.Scheme, &sch); err != nil {
		cancel()
		return &formsStore.CreateOneResponse{}, nil
	}
	if err := json.Unmarshal(request.UiScheme, &uiSch); err != nil {
		cancel()
		return &formsStore.CreateOneResponse{}, nil
	}
	model.FormName = request.Name
	model.Scheme = sch
	model.UiScheme = uiSch

	result, err := s.collection.InsertOne(ctx, model)
	if err != nil {
		cancel()
		return &formsStore.CreateOneResponse{}, nil
	}
	_, err = s.collection.InsertOne(ctx, model)
	if err != nil {
		cancel()
		return &formsStore.CreateOneResponse{}, nil
	}
	cancel()
	return &formsStore.CreateOneResponse{Id: cast.ToString(result.InsertedID)}, nil
}

func (s *store) GetOne(request formsStore.GetOneRequest) (*formsStore.GetOneResponse, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
	var formsModel forms.Forms

	objectId, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil{
		fmt.Println("Invalid id")
		cancel()
		return &formsStore.GetOneResponse{}, err
	}

	err = s.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&formsModel)
	if err != nil {
		fmt.Println("No forms found")
		cancel()
		return &formsStore.GetOneResponse{}, err
	}
	cancel()
	return &formsStore.GetOneResponse{
		Name: formsModel.FormName,
		Scheme: formsModel.Scheme,
		UiScheme: formsModel.UiScheme,
	}, nil
}

