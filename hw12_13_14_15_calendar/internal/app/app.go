package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	storage "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage Storage
	logger  *slog.Logger
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	ListEvents(ctx context.Context) ([]storage.Event, error)
	DeleteEvent(ctx context.Context, id string) error
	UpdateEvent(ctx context.Context, newEvent storage.Event) error
}

func New(log *slog.Logger, storage Storage) *App {
	return &App{
		storage: storage,
		logger:  log,
	}
}

func (a *App) CreateEvent(
	ctx context.Context,
	id, title *string,
	date *time.Time,
	duration *string,
	descr *string,
	userID *string,
	notificationTime *string,
) error {
	a.logger.Info("CreateEvent", "id", *id, "title", *title)
	return a.storage.CreateEvent(
		ctx,
		storage.Event{
			ID:               id,
			Title:            title,
			Date:             date,
			Duration:         duration,
			Description:      descr,
			UserID:           userID,
			NotificationTime: notificationTime,
		},
	)
}

func (a *App) ListEvents(ctx context.Context) (events []storage.Event, err error) {
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
	event storage.Event,
) error {
	a.logger.Info("UpdateEvent", "id", *event.ID)
	return a.storage.UpdateEvent(
		ctx,
		event,
	)
}
