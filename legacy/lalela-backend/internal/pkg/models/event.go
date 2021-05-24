package models


//// EventCon

type EventLogRequest struct {
	Event Event
}

type EventLogResponse struct {
	Messages []string `json:"message"`
}


