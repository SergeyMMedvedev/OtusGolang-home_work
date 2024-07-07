package memorystorage

import (
	"context"
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
