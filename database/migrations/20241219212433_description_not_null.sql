-- +goose Up
-- +goose StatementBegin
-- First ensure no NULL values exist
UPDATE events SET description = '' WHERE description IS NULL;

-- Create new table with desired schema
CREATE TABLE new_events
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    title          VARCHAR(255) NOT NULL,
    description    TEXT NOT NULL DEFAULT '',
    event_date     DATETIME     NOT NULL,
    location       VARCHAR(255),
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Copy data from old table to new table
INSERT INTO new_events SELECT * FROM events;

-- Drop old table
DROP TABLE events;

-- Rename new table to original name
ALTER TABLE new_events RENAME TO events;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert to original schema
CREATE TABLE new_events
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    title          VARCHAR(255) NOT NULL,
    description    TEXT,
    event_date     DATETIME     NOT NULL,
    location       VARCHAR(255),
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Copy data
INSERT INTO new_events SELECT * FROM events;

-- Drop modified table
DROP TABLE events;

-- Rename back to original
ALTER TABLE new_events RENAME TO events;
-- +goose StatementEnd