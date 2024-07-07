package storage

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
	switch config.Type {
	case "psql":
		storage = sqlstorage.New(config.Psql)
	case "memory":
		storage = memorystorage.New()
	default:
		return nil, errors.New("unknown storage type")
	}
	return storage, nil
}
