-- +goose Up
-- +goose StatementBegin
CREATE TABLE events
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    title          VARCHAR(255) NOT NULL,
    description    TEXT,
    event_date     DATETIME     NOT NULL,
    location       VARCHAR(255),
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
