package sqlstorage

import (
	"context"
	"fmt"
	"reflect"

	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for sqlx
	"github.com/pressly/goose/v3"
)

type Storage struct { // TODO
	Host      string
	Port      int64
	User      string
	Pass      string
	DBName    string
	Sslmode   string
	Migration string

	db *sqlx.DB
}

func New(conf c.PsqlConf) *Storage {
	return &Storage{
		Host:      conf.Host,
		Port:      conf.Port,
		User:      conf.User,
		Pass:      conf.Password,
		DBName:    conf.Dbname,
		Sslmode:   conf.Sslmode,
		Migration: conf.Migration,
	}
}

func (s *Storage) getDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.Host, s.Port, s.User, s.Pass, s.DBName, s.Sslmode,
	)
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.db, err = sqlx.Open("postgres", s.getDSN())
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	return s.db.PingContext(ctx)
}

func (s *Storage) Migrate(_ context.Context) (err error) {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}
	if err := goose.Up(s.db.DB, s.Migration); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) ListEvents(ctx context.Context) (events []storage.Event, err error) {
	query := "select * from events"
	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	for rows.Next() {
		var event storage.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	fmt.Printf("%+v\n", event)
	query := `
	insert into events 
	(id, title, date, duration, description, user_id, notification_time) 
	values 
	($1, $2, $3, $4, $5, $6, $7)`
	result, err := s.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.Title,
		event.Date,
		event.Duration,
		event.Description,
		event.UserID,
		event.NotificationTime,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affectedRows == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	query := "delete from events where id = $1"
	result, err := s.db.ExecContext(ctx, query, eventID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affectedRows == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, newEvent storage.Event) error {
	query := "UPDATE events SET "

	args := []interface{}{}
	idx := 1
	newEventValue := reflect.ValueOf(newEvent)
	for i := 0; i < newEventValue.NumField(); i++ {
		v := newEventValue.Field(i)
		if !v.IsNil() {
			query += fmt.Sprintf("%s = $%d, ", newEventValue.Type().Field(i).Name, idx)
			args = append(args, v.Interface())
			idx++
		}
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d", idx)
	args = append(args, newEvent.ID)
	result, err := s.db.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affectedRows == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}
