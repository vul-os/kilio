package store

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lalela-backend/internal/pkg/database"
	"lalela-backend/internal/pkg/submissions"
	"time"
)

func CreateOne(formId string, orgId string, submissionData json.RawMessage, identifier string) (error){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var model submissions.Submissions
	var collection = database.OpenCollection("submissions")

	var submish interface{}
	if err := json.Unmarshal(submissionData, &submish); err != nil {
		cancel()
		return err
	}
	cancel()

	model.FormId, _ = primitive.ObjectIDFromHex(formId)
	model.OrganizationId, _ = primitive.ObjectIDFromHex(orgId)
	model.SubmissionData = submissionData
	model.Identifier = identifier
	model.Status = "submitted"

	_, err := collection.InsertOne(ctx, model)
	if err != nil {
		cancel()
		return err
	}
	cancel()
	return nil
}
