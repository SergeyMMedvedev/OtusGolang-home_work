package main

import (
	"context"
	"errors"

	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

type Storage interface { // TODO
	app.Storage

	Connect(ctx context.Context) error
	Migrate(ctx context.Context) error
}

func NewStorage(config config.StorageConf) (storage Storage, err error) {
	if config.Type == "psql" {
		storage = sqlstorage.New(config.Psql)
	} else if config.Type == "memory" {
		storage = memorystorage.New()
	} else {
		return nil, errors.New("Unknown storage type")
	}
	return storage, nil
}
