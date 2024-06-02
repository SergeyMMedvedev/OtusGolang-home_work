package memorystorage

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	events := make(map[string]storage.Event)
	return &Storage{
		events: events,
	}
}

func (s *Storage) ListEvents(_ context.Context) (events []storage.Event, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		events = append(events, event)
	}

	return
}

func (s *Storage) CreateEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[*event.ID] = event

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, eventID)

	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, newEvent storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	event := s.events[*newEvent.ID]
	newEventValue := reflect.ValueOf(newEvent)
	oldEventValue := reflect.ValueOf(&event).Elem()
	for i := 0; i < newEventValue.NumField(); i++ {
		newField := newEventValue.Field(i)
		oldField := oldEventValue.Field(i)
		if oldField.CanSet() {
			if !newField.IsNil() {
				oldField.Elem().Set(newField.Elem())
			}
		} else {
			return fmt.Errorf("field %s is not settable", oldField.Type().Name())
		}
	}
	return nil
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Migrate(_ context.Context) error {
	return nil
}

// TODO
