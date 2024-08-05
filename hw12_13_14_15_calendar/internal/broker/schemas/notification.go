package schemas

import (
	"fmt"
	"time"
)

type Notification struct {
	EventID    string
	EventTitle string
	EventDate  time.Time
	UserID     string
}

func (n Notification) String() string {
	return fmt.Sprintf(
		"{\"EventId\": \"%s\", \"EventTitle\": \"%s\", \"EventDate\": \"%s\", \"UserID\": \"%s\"}",
		n.EventID, n.EventTitle, n.EventDate.Format("2006-01-02T15:04:05Z07:00"), n.UserID,
	)
}
