package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/services"
	"net/http"
	"time"
)

type FormsCon struct{}

type FormCreateRequest struct {
	Name string `json:"name"`
	Scheme json.RawMessage `json:"scheme"`
	UiScheme json.RawMessage `json:"ui_scheme"`
}

type FormCreateResponse struct {
	Id string `json:"id"`
}


func (t *FormsCon) CreateForm(r *http.Request, args *FormCreateRequest,	reply *FormCreateResponse) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
	var model models.Forms
	var collection = services.OpenCollection("forms")

	var sceheme interface{}
	var uiSceheme interface{}

	if err := json.Unmarshal(args.Scheme, &sceheme); err != nil {
		cancel()
		return err
	}
	if err := json.Unmarshal(args.UiScheme, &uiSceheme); err != nil {
		cancel()
		return err
	}
	model.FormName = args.Name
	model.Scheme = sceheme
	model.UiScheme = uiSceheme

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


type FormGetRequest struct {
	Id string `json:"id"`
}


func (t *FormsCon) GetForm(r *http.Request, args *FormGetRequest, reply *FormCreateRequest) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
	var formsModel models.Forms

	var collection = services.OpenCollection("forms")

	objectId, err := primitive.ObjectIDFromHex(args.Id)
	if err != nil{
		fmt.Println("Invalid id")
		cancel()
		return err
	}

	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&formsModel)
	if err != nil {
		fmt.Println("No forms found")
		cancel()
		return err
	}
	defer cancel()

	fmt.Println(formsModel)
	cancel()
	return nil
}



