package models

import (
	"database/sql"
	"errors"
	"time"
)

// Event represents a scheduled occurrence with a title, description, date, and location.
type Event struct {
	Id          int
	Title       string
	Description string
	EventDate   time.Time
	Location    string
}

// EventModel provides methods for managing and interacting with events in the database.
// It includes functionality to create, retrieve, and list event records.
// The `DB` field holds the database connection used for queries and operations.
type EventModel struct {
	DB *sql.DB
}

// Create adds a new event record to the database with the provided title, description, date, and location.
// It returns the ID of the newly created event or an error if the operation fails.
func (m *EventModel) Create(title, description string, eventDate time.Time, location string) (int, error) {
	return 0, nil
}

// Get retrieves an event from the database by its unique ID.
// It returns the matching Event object or an error if the query fails or no event is found.
func (m *EventModel) Get(id int) (Event, error) {

	stmt := "SELECT * FROM events WHERE id = ?"

	row := m.DB.QueryRow(stmt, id)

	var e Event

	err := row.Scan(&e.Id, &e.Title, &e.Description, &e.EventDate, &e.Location)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return Event{}, ErrNoRecord
		} else {
			return Event{}, err
		}
	}
	return e, nil
}

// GetAll retrieves all event records from the database and returns them as a slice of Event pointers or an error if it fails.
func (m *EventModel) GetAll() ([]*Event, error) {
	return nil, nil
}
