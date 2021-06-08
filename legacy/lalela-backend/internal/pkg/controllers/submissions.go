package controllers

import (
	"encoding/json"
	"lalela-backend/internal/pkg/services"
	"net/http"
)

type SubmissionsCon struct{}

type SubmissionSubmitRequest struct {
	FormId         string          `json:"form_id"`
	OrganizationId string          `json:"organization_id"`
	SubmissionData json.RawMessage `json:"submission_data"`
	Identifier     string          `json:"identifier"`
}

type SubmissionSubmitResponse struct {
	Id         string `json:"id"`
	Identifier string `json:"identifier"`
}

func (t *SubmissionsCon) SubmitForm(r *http.Request, args *SubmissionSubmitRequest, reply *SubmissionSubmitResponse) error {
	err := services.AddForm(args.FormId, args.OrganizationId, args.SubmissionData, args.Identifier)
	if err != nil {
		return err
	}
	return nil
}

type GetSubmissionsInOrgRequest struct {
	OrganizationId string `json:"identifier"`
}

func (t *SubmissionsCon) GetSubmissionsInOrg(r *http.Request, args *SubmissionSubmitRequest, reply *SubmissionSubmitResponse) error {
	canDo, err := services.Authorize(r, args.OrganizationId, "submissions", "get")
	if err != nil {
		return err
	}
	if !canDo {
		return nil
	}

	return nil
}
