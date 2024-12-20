package main

import (
	"github.com/madalinpopa/go-event-planner/internal/validator"
	"time"
)

type EventForm struct {
	Title               string    `form:"title"`
	Description         string    `form:"description"`
	Location            string    `form:"location"`
	EventDate           time.Time `form:"eventDate"`
	validator.Validator `form:"-"`
}
