package main

import (
	"github.com/madalinpopa/go-event-planner/internal/validator"
	"time"
)

// EventForm represents a data structure for handling event creation form input and validation.
type EventForm struct {
	Title               string    `form:"title"`
	Description         string    `form:"description"`
	Location            string    `form:"location"`
	EventDate           time.Time `form:"eventDate"`
	validator.Validator `form:"-"`
}

// UserRegisterForm represents the structure for a user registration form containing name, email, and password fields.
type UserRegisterForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
