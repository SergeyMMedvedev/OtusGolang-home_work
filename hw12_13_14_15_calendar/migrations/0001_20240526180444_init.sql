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
    notification_time INTERVAL
);
-- CREATE INDEX idx_event_user_id ON event(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- drop index idx_event_user_id;
drop table events;
-- +goose StatementEnd


-- CREATE TABLE event (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),    -- Уникальный идентификатор события
--     title VARCHAR(255) NOT NULL,                      -- Заголовок (короткий текст)
--     date TIMESTAMP NOT NULL,                      -- Дата и время события
--     duration INTERVAL,                                -- Длительность события (или дата и время окончания)
--     description TEXT,                                 -- Описание события (длинный текст, опционально)
--     user_id UUID NOT NULL,                            -- ID пользователя, владельца события
--     notification_time INTERVAL                        -- За сколько времени высылать уведомление (опционально)
-- );