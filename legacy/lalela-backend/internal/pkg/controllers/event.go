package controllers

import "cog-analytics-engine-go/internal/pkg/models"

//// EventCon

type EventLogRequest struct {
	Event models.Event
}

type EventLogResponse struct {
	Messages []string `json:"message"`
}


