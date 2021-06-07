package controllers

import (
	"context"
	"encoding/json"
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/services"
	"net/http"
	"time"
)

type SubmissionsCon struct{}

type SubmissionSubmitRequest struct {
	FormId         string          `json:"form_id"`
	SubmissionData json.RawMessage `json:"submission_data"`
	Identifier     string          `json:"identifier"`
}

type SubmissionSubmitResponse struct {
	Id         string `json:"id"`
	Identifier string `json:"identifier"`
}

func (t *SubmissionsCon) SubmitForm(r *http.Request, args *SubmissionSubmitRequest, reply *SubmissionSubmitResponse) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var model models.Submissions
	var collection = services.OpenCollection("submissions")

	var submish interface{}
	if err := json.Unmarshal(args.SubmissionData, &submish); err != nil {
		cancel()
		return err
	}

	model.FormId = args.FormId
	model.SubmissionData = args.SubmissionData
	model.Identifier = args.Identifier
	model.Status = "submitted"

	_, err := collection.InsertOne(ctx, model)
	if err != nil {
		cancel()
		return err
	}
	cancel()
	return nil
}
