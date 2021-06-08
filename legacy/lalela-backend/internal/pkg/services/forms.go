package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lalela-backend/internal/pkg/models"
	"time"
)

func CreateForm(name string, scheme json.RawMessage, uischeme json.RawMessage) (error){
	var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
	var model models.Forms
	var collection = OpenCollection("forms")

	var sch interface{}
	var uiSch interface{}

	if err := json.Unmarshal(scheme, &sch); err != nil {
		cancel()
		return err
	}
	if err := json.Unmarshal(uischeme, &uiSch); err != nil {
		cancel()
		return err
	}
	model.FormName = name
	model.Scheme = sch
	model.UiScheme = uiSch

	_, err := collection.InsertOne(ctx, model)
	if err != nil {
		cancel()
		return err
	}
	_, err = collection.InsertOne(ctx, model)
	if err != nil {
		cancel()
		return err
	}
	cancel()
	return nil
}

func GetForm(id string) (models.Forms, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
	var formsModel models.Forms

	var collection = OpenCollection("forms")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		fmt.Println("Invalid id")
		cancel()
		return models.Forms{}, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&formsModel)
	if err != nil {
		fmt.Println("No forms found")
		cancel()
		return models.Forms{}, err
	}
	cancel()
	return formsModel, nil
}
