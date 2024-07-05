package schemas

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Event struct {
	ID               string
	Title            string
	Date             time.Time
	Duration         string `default:"00:30:00"`
	Description      string
	UserID           string `db:"user_id"`
	NotificationTime int32  `db:"notification_time"`
}

func (e Event) parseToDurationPb(s string) (*durationpb.Duration, error) {
	if strings.Contains(s, ":") {
		var hour, min, sec int
		_, err := fmt.Sscanf(s, "%d:%d:%d", &hour, &min, &sec)
		if err != nil {
			slog.Error(fmt.Sprintf("Error parsing duration: %s", err.Error()))
			return nil, err
		}
		duration := time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
		return durationpb.New(duration), nil
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		slog.Error(fmt.Sprintf("Error parsing duration: %s", err.Error()))
		return nil, err
	}
	return durationpb.New(duration), nil
}

func (e Event) GetDurationPb() (*durationpb.Duration, error) {
	return e.parseToDurationPb(e.Duration)
}

func (e Event) String() string {
	return fmt.Sprintf(
		"{ID: %s, Title: %s, Date: %s, Duration: %s, Description: %s, UserID: %s, NotificationTime: %v}",
		e.ID, e.Title, e.Date, e.Duration, e.Description, e.UserID, e.NotificationTime,
	)
}

func NewEvent() Event {
	newID := uuid.New().String()
	return Event{
		ID:               newID,
		NotificationTime: 1,
	}
}
