package models

import (
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

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
func (m *UserModel) Create(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"

	_, err = m.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		var sqliteError *sqlite3.Error
		if errors.As(err, &sqliteError) && errors.Is(sqliteError.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

// Authenticate verifies a user's credentials and returns
// their ID if valid or an error if authentication fails.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

// Exists checks if a user with the specified ID exists in the
// database and returns a boolean result and an error.
func (m *UserModel) Exists(id int) (bool, error) {

	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
