package controllers

import "cog-analytics-engine-go/internal/pkg/models"

type DateProcessedRequest struct {
	Email string `json:"email"`
}

type DateProcessedResponse struct {
	Messages []string `json:"message"`
}

type DateProcessedLogRequest struct {
	Email string `json:"email"`
}

type DateProcessedLogResponse struct {
	Logs []models.DateProcessedItem `json:"log"`
}

type DateProcessedTableResponse struct {
	Logs map[string]string `json:"log"`
}
