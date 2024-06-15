package app

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"

	schemas "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage/schemas"
)

type App struct {
	storage Storage
	logger  *slog.Logger
}

type Storage interface {
	CreateEvent(ctx context.Context, event schemas.Event) error
	ListEvents(ctx context.Context) ([]schemas.Event, error)
	DeleteEvent(ctx context.Context, id string) error
	UpdateEvent(ctx context.Context, newEvent schemas.Event) error
}

func New(log *slog.Logger, storage Storage) *App {
	return &App{
		storage: storage,
		logger:  log,
	}
}

func (a *App) CreateEvent(
	ctx context.Context,
	event schemas.Event,
) error {
	a.logger.Info("CreateEvent", "id", event.ID, "title", event.Title)
	return a.storage.CreateEvent(
		ctx,
		schemas.Event{
			ID:               uuid.New().String(),
			Title:            event.Title,
			Date:             event.Date,
			Duration:         event.Duration,
			Description:      event.Description,
			UserID:           event.UserID,
			NotificationTime: event.NotificationTime,
		},
	)
}

func (a *App) ListEvents(ctx context.Context) (events []schemas.Event, err error) {
	events, err = a.storage.ListEvents(ctx)
	if err != nil {
		a.logger.Error("ListEvents", "err", err)
		return nil, fmt.Errorf("app ListEvents error: %w", err)
	}
	return events, nil
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	a.logger.Info("DeleteEvent", "id", id)
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) UpdateEvent(
	ctx context.Context,
	event schemas.Event,
) error {
	a.logger.Info("UpdateEvent", "id", event.ID)
	return a.storage.UpdateEvent(
		ctx,
		event,
	)
}
