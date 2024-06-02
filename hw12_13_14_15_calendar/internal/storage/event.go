package storage

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               *string
	Title            *string
	Date             *time.Time
	Duration         *string `default:"00:30:00"`
	Description      *string
	UserID           *string `db:"user_id"`
	NotificationTime *string `db:"notification_time"`
}

func (e Event) String() string {
	return fmt.Sprintf(
		"{ID: %s, Title: %s, Date: %s, Duration: %s, Description: %s, UserID: %s, NotificationTime: %s}",
		*e.ID, *e.Title, *e.Date, *e.Duration, *e.Description, *e.UserID, *e.NotificationTime,
	)
}

func NewEvent() Event {
	defaultTime := "00:30:00"
	newID := uuid.New().String()
	return Event{
		ID:               &newID,
		NotificationTime: &defaultTime,
	}
}
