package grpcserver

import (
	"time"
)

func StartOfWeek(date time.Time) time.Time {
	weekday := date.Weekday()
	var daysToSubtract int
	if weekday == time.Sunday {
		daysToSubtract = 6
	} else {
		daysToSubtract = int(weekday) - int(time.Monday)
	}
	startOfWeek := date.AddDate(0, 0, -daysToSubtract)
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	return startOfWeek
}
