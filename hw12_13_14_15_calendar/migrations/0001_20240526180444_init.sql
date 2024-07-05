-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    duration INTERVAL,
    description TEXT DEFAULT '',
    user_id UUID NOT NULL,
    notification_time INT
);
-- CREATE INDEX idx_event_user_id ON event(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- drop index idx_event_user_id;
drop table events;
-- +goose StatementEnd


