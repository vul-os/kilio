package models

import "time"

// Overview

type OverviewEventsLog struct {
	EventID        uint
	EventUser      string
	EventUserGroup string
	EventType      string
	EventDate      time.Time
}

type OverviewGetRequest struct {
	Email string `json:"email"`
}

type OverviewGetResponse struct {
	Users         int64
	Dashboards    int64
	Groups        int64
	Logins        int64
	Events        []OverviewEventsLog
	DateProcessed []DateProcessedItem
}
