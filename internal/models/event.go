package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

// Event represents a scheduled occurrence with a title, description, date, and location.
type Event struct {
	Id          int
	Title       string
	Description string
	Location    string
	EventDate   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

	stmt := `INSERT INTO events (title, description, event_date, location) 
	VALUES (?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, title, description, eventDate, location)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Delete removes an event record from the database by its unique ID and returns an error if the operation fails.
func (m *EventModel) Delete(id int) error {
	stmt := "DELETE FROM events WHERE id = ?"

	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to fetch rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no event found with id: %d", id)
	}

	return nil
}

// Get retrieves an event from the database by its unique ID.
// It returns the matching Event object or an error if the query fails or no event is found.
func (m *EventModel) Get(id int) (Event, error) {

	stmt := "SELECT * FROM events WHERE id = ?"

	row := m.DB.QueryRow(stmt, id)

	var e Event

	err := row.Scan(&e.Id, &e.Title, &e.Description, &e.EventDate, &e.Location, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Event{}, ErrNoRecord
		} else {
			return Event{}, err
		}
	}
	return e, nil
}

// GetAll retrieves all event records from the database and returns them as a slice of Event pointers or an error if it fails.
func (m *EventModel) GetAll() ([]Event, error) {
	stmt := "SELECT id, title, description, event_date, location, created_at, updated_at FROM events"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}(rows)

	var events []Event
	for rows.Next() {

		var e Event
		err := rows.Scan(&e.Id, &e.Title, &e.Description, &e.EventDate, &e.Location, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
