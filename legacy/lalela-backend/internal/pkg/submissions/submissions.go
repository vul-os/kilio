package submissions

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Submissions struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FormId         primitive.ObjectID `json:"form_id"`
	OrganizationId primitive.ObjectID `json:"organization_id"`
	SubmissionData interface{}        `json:"submission_data"`
	Identifier     string             `json:"identifier"`
	Status         string             `json:"status"`
}

