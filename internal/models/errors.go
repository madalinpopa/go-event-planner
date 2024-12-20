package models

import "errors"

var (

	// ErrNoRecord indicates that no matching record was found in the database or data source.
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrInvalidCredentials indicates that the provided user credentials are
	// invalid or do not match any record in the system.
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// ErrDuplicateEmail indicates that the provided email already exists
	// in the system and cannot be used again.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
