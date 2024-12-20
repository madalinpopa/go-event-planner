package models

import "database/sql"

// User represents a user entity with basic identification and authentication fields.
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

// UserModel provides methods to interact with the users' data in the database using an sql.DB instance.
type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Create(name, email, password string) error {
	_ = name
	_ = email
	_ = password

	return nil
}

func (m *UserModel) Get(id int) (User, error) {
	_ = id
	return User{}, nil
}
