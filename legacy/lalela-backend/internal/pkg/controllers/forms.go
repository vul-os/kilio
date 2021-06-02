package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/utils"
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
	var formsModel models.Forms

	var formsCollection = utils.OpenCollection(utils.Client, "forms")

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
	formsModel.FormName = args.Name
	formsModel.Scheme = sceheme
	formsModel.UiScheme = uiSceheme

	_, err := formsCollection.InsertOne(ctx, formsModel)
	if err != nil {
		cancel()
		return err
	}
	_, err = formsCollection.InsertOne(ctx, formsModel)
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

	var formsCollection = utils.OpenCollection(utils.Client, "forms")

	objectId, err := primitive.ObjectIDFromHex("60b8080ddf18fdc3625491db")
	if err != nil{
		fmt.Println("Invalid id")
		cancel()
		return err
	}

	err = formsCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&formsModel)
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

