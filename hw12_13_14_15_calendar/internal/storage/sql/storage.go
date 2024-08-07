package sqlstorage

import (
	"context"
	"fmt"
	"time"

	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage/schemas"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for sqlx
	"github.com/pressly/goose/v3"
)

type Storage struct { // TODO
	Host          string
	Port          int64
	User          string
	Pass          string
	DBName        string
	Sslmode       string
	MigrationDir  string
	ExecMigration bool

	db *sqlx.DB
}

func New(conf c.PsqlConf) *Storage {
	return &Storage{
		Host:          conf.Host,
		Port:          conf.Port,
		User:          conf.User,
		Pass:          conf.Password,
		DBName:        conf.Dbname,
		Sslmode:       conf.Sslmode,
		MigrationDir:  conf.MigrationDir,
		ExecMigration: conf.ExecMigration,
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
	if s.ExecMigration {
		if err := goose.SetDialect("postgres"); err != nil {
			return fmt.Errorf("cannot set dialect: %w", err)
		}
		if err := goose.Up(s.db.DB, s.MigrationDir); err != nil {
			return fmt.Errorf("cannot do up migration: %w", err)
		}
	}
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) ListEvents(ctx context.Context) (events []schemas.Event, err error) {
	query := "select * from events"
	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	for rows.Next() {
		var event schemas.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Storage) ListDayEvents(ctx context.Context, date time.Time) (events []schemas.Event, err error) {
	// select from events where date day, month and year equal event date day, month and year
	query := `
	select * from events
	where
	date_part('day', date) = date_part('day', $1::timestamp)
	and date_part('month', date) = date_part('month', $1::timestamp)
	and date_part('year', date) = date_part('year', $1::timestamp)
	`
	rows, err := s.db.QueryxContext(ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	for rows.Next() {
		var event schemas.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}
	return events, nil
}

func (s *Storage) ListEventsForNotification(ctx context.Context) (events []schemas.Event, err error) {
	query := `
	select * from events
	where
	DATE_TRUNC('second', current_timestamp AT TIME ZONE $1) = 
	DATE_TRUNC('second', date) - (notification_time || ' days')::INTERVAL;
	`
	zoneName, _ := time.Now().Zone()
	rows, err := s.db.QueryxContext(ctx, query, zoneName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	for rows.Next() {
		var event schemas.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}
	return events, nil
}

func (s *Storage) ListWeekEvents(ctx context.Context, date time.Time) (events []schemas.Event, err error) {
	query := `
	select * from events
	where
	date_part('week', date) = date_part('week', $1::timestamp)
	and date_part('month', date) = date_part('month', $1::timestamp)
	and date_part('year', date) = date_part('year', $1::timestamp)
	`
	rows, err := s.db.QueryxContext(ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	for rows.Next() {
		var event schemas.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Storage) ListMonthEvents(ctx context.Context, date time.Time) (events []schemas.Event, err error) {
	query := `
	select * from events
	where
	date_part('month', date) = date_part('month', $1::timestamp)
	and date_part('year', date) = date_part('year', $1::timestamp)
	`
	rows, err := s.db.QueryxContext(ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	for rows.Next() {
		var event schemas.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Storage) ListLastYearEvents(ctx context.Context) (events []schemas.Event, err error) {
	query := `
	select * from events
	where
	DATE_TRUNC('second', date) < (DATE_TRUNC('second', current_timestamp AT TIME ZONE $1) -(1 || ' years')::INTERVAL);
	`
	zoneName, _ := time.Now().Zone()
	rows, err := s.db.QueryxContext(ctx, query, zoneName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	for rows.Next() {
		var event schemas.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Storage) CreateEvent(ctx context.Context, event schemas.Event) error {
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

func (s *Storage) UpdateEvent(ctx context.Context, newEvent schemas.Event) error {
	query := `
	UPDATE events 
	SET 
	title = :title,
	date = :date,
	duration = :duration,
	description = :description,
	user_id = :user_id,
	notification_time = :notification_time
	where id = :id
	`
	result, err := s.db.NamedExecContext(
		ctx,
		query,
		newEvent,
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
