package adapter

import (
	"encoding/json"
	"lalela-backend/internal/pkg/submissions/store"
	"net/http"
)

type SubmissionsCon struct{}


type CreateOneRequest struct {
	FormId         string          `json:"form_id"`
	OrganizationId string          `json:"organization_id"`
	SubmissionData json.RawMessage `json:"submission_data"`
	Identifier     string          `json:"identifier"`
}

type CreateOneResponse struct {
	Id         string `json:"id"`
	Identifier string `json:"identifier"`
}

func (t *SubmissionsCon) CreateOne(r *http.Request, args *CreateOneRequest, reply *CreateOneResponse) error {
	err := store.CreateOne(args.FormId, args.OrganizationId, args.SubmissionData, args.Identifier)
	if err != nil {
		return err
	}
	return nil
}

//type GetSubmissionsInOrgRequest struct {
//	OrganizationId string `json:"identifier"`
//}
//
//func (t *SubmissionsCon) GetOne(r *http.Request, args *SubmissionSubmitRequest, reply *SubmissionSubmitResponse) error {
//	canDo, err := auth.Authorize(r, "submissions", "get")
//	if err != nil {
//		return err
//	}
//	if !canDo {
//		return nil
//	}
//
//	return nil
//}

