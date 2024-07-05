package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage/schemas"
)

type Storage struct {
	events map[string]schemas.Event
	mu     sync.RWMutex
}

func New() *Storage {
	events := make(map[string]schemas.Event)
	return &Storage{
		events: events,
	}
}

func (s *Storage) ListEvents(_ context.Context) (events []schemas.Event, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		events = append(events, event)
	}

	return
}

func (s *Storage) ListEventsForNotification(_ context.Context) (events []schemas.Event, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		t := event.Date.Add(-time.Duration(event.NotificationTime) * time.Hour).Truncate(time.Second)
		now := time.Now().Truncate(time.Second)
		if t == now {
			events = append(events, event)
		}
	}

	return
}

func (s *Storage) ListDayEvents(_ context.Context, date time.Time) (events []schemas.Event, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		if event.Date.Year() == date.Year() && event.Date.Month() == date.Month() && event.Date.Day() == date.Day() {
			events = append(events, event)
		}
		events = append(events, event)
	}
	return
}

func (s *Storage) ListWeekEvents(_ context.Context, date time.Time) (events []schemas.Event, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		if event.Date.Year() == date.Year() &&
			event.Date.Month() == date.Month() &&
			event.Date.Day() < date.Day()+7 &&
			event.Date.Day() >= date.Day() {
			events = append(events, event)
		}
	}
	return
}

func (s *Storage) ListMonthEvents(_ context.Context, date time.Time) (events []schemas.Event, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		if event.Date.Year() == date.Year() && event.Date.Month() == date.Month() {
			events = append(events, event)
		}
	}
	return
}

func (s *Storage) CreateEvent(_ context.Context, event schemas.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, eventID)

	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, newEvent schemas.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[newEvent.ID]; ok {
		s.events[newEvent.ID] = newEvent
		return nil
	}
	return fmt.Errorf("event with id %s not found", newEvent.ID)
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Migrate(_ context.Context) error {
	return nil
}

// TODO
