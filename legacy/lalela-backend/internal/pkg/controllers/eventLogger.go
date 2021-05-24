package controllers

import (
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/utils"
	"time"
)

func EventLog(UserID string, Event string, Object string, Type string) error {
	db := utils.GetDB()
	event := &models.Event{
		UserID: UserID,
		Event:  Event,
		Object: Object,
		Type:   Type,
		Date:   time.Now(),
	}
	// Log Event
	db.Create(event)
	return nil
}
