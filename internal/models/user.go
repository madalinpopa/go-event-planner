package models

import "database/sql"

// User represents a user entity with basic
// identification and authentication fields.
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
}

// UserModel provides methods to interact with the users'
// data in the database using an sql.DB instance.
type UserModel struct {
	DB *sql.DB
}

// Create adds a new user with the provided name,
// email, and hashed password to the database.
func (m *UserModel) Create(name, email string, password []byte) error {
	return nil
}

// Authenticate verifies a user's credentials and returns
// their ID if valid or an error if authentication fails.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Exists checks if a user with the specified ID exists in the
// database and returns a boolean result and an error.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
