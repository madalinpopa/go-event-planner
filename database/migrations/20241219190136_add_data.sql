-- +goose Up
-- +goose StatementBegin
INSERT INTO events (title, description, event_date, location)
VALUES
    /* Format: YYYY-MM-DD HH:MM:SS */
    ('Tech Conference 2025',
     'Annual technology conference featuring the latest innovations in AI and cloud computing',
     '2025-03-15 09:00:00',
     'Convention Center, San Francisco'),

    ('Team Building Workshop',
     'Interactive workshop focusing on leadership and collaboration skills',
     datetime('2024-12-28 13:00:00'),
     'Downtown Business Center'),

    ('Product Launch Event',
     'Launch of our new software platform with live demonstrations',
     '2025-01-20 18:30:00',
     'Tech Hub, Austin'),

    ('Agile Development Seminar',
     'Learn best practices in agile methodology and scrum framework',
     '2025-02-10 10:00:00',
     'Online - Virtual Event'),

    ('End of Year Party',
     'Celebrate our achievements and success throughout the year',
     '2024-12-31 19:00:00',
     'Grand Hotel Ballroom');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM events
WHERE title IN (
                'Tech Conference 2025',
                'Team Building Workshop',
                'Product Launch Event',
                'Agile Development Seminar',
                'End of Year Party'
    );
-- +goose StatementEnd