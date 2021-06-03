package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Submissions struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FormId         string             `json:"form_id"`
	SubmissionData interface{}        `json:"submission_data"`
	Identifier     string             `json:"identifier"`
	Status         string			  `json:"status"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}
