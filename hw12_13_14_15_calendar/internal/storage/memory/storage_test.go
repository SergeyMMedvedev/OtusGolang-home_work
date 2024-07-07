package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	schemas "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage/schemas"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func getNewEvent() schemas.Event {
	event := schemas.NewEvent()

	event.ID = uuid.New().String()
	event.Title = "Event title"
	event.Date = time.Now().Add(time.Hour * 24)
	event.Duration = "30:00:00"
	event.Description = "Event description"
	event.UserID = uuid.New().String()
	return event
}

func getNewEventWithParams(date time.Time, ntTime int32) schemas.Event {
	event := schemas.NewEvent()

	event.ID = uuid.New().String()
	event.Title = "Event title"
	event.Date = date
	event.Duration = "30:00:00"
	event.Description = "Event description"
	event.UserID = uuid.New().String()
	event.NotificationTime = ntTime
	return event
}

func TestStorage(t *testing.T) {
	s := New()
	ctx := context.Background()
	event := getNewEvent()

	events, err := s.ListEvents(ctx)
	require.NoError(t, err)
	require.Empty(t, events)

	err = s.CreateEvent(ctx, event)
	require.NoError(t, err)

	events, err = s.ListEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, 1)

	newTitle := "New title"
	err = s.UpdateEvent(ctx, schemas.Event{
		ID:    event.ID,
		Title: newTitle,
	})
	require.NoError(t, err)

	events, err = s.ListEvents(ctx)
	event = events[0]
	require.NoError(t, err)
	require.Equal(t, newTitle, event.Title)

	err = s.DeleteEvent(ctx, event.ID)
	require.NoError(t, err)
	events, err = s.ListEvents(ctx)
	require.NoError(t, err)
	require.Empty(t, events)
}

func TestStorageConcurrencyCreation(t *testing.T) {
	s := New()
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			event := getNewEvent()
			defer wg.Done()
			err := s.CreateEvent(ctx, event)
			require.NoError(t, err)
		}()
	}

	wg.Wait()

	events, err := s.ListEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, 1000)
}

func TestListEventsForNotification(t *testing.T) {
	s := New()
	ctx := context.Background()
	now := time.Now()

	d := now.Add(time.Hour*24 + time.Second*3)
	ntTime := int32(1)
	event := getNewEventWithParams(d, ntTime)

	events, err := s.ListEvents(ctx)
	require.NoError(t, err)
	require.Empty(t, events)

	err = s.CreateEvent(ctx, event)
	require.NoError(t, err)

	ticker := time.NewTicker(time.Second)
	done := make(chan bool, 1)
	stopTime := now.Add(time.Second * 10)
	var eventsForNt []schemas.Event
OUTER:
	for {
		select {
		case <-done:
			fmt.Println("Done")
			ticker.Stop()
			break OUTER
		case <-ticker.C:
			now := time.Now()
			fmt.Println("Check List Events For Notification " + now.String())
			if now.After(stopTime) {
				fmt.Println("Events for notification not founded!")
				done <- true
			}
			eventsForNt, err = s.ListEventsForNotification(ctx)
			require.NoError(t, err)
			if len(eventsForNt) > 0 {
				fmt.Println("Events for notification founded!")
				done <- true
			}
		}
	}
	for _, e := range eventsForNt {
		fmt.Println(e)
	}
	require.NoError(t, err)
	require.Len(t, eventsForNt, 1)
	event = eventsForNt[0]
	err = s.DeleteEvent(ctx, event.ID)
	require.NoError(t, err)
	events, err = s.ListEvents(ctx)
	require.NoError(t, err)
	require.Empty(t, events)
}
