package main

import (
	"errors"
	"fmt"
	"github.com/madalinpopa/go-event-planner/internal/models"
	"github.com/madalinpopa/go-event-planner/internal/validator"
	"net/http"
	"runtime/debug"
	"strconv"
)

// ping handles the /ping endpoint, responding with "pong" to indicate the service is available and operational.
func (app *App) ping(w http.ResponseWriter, r *http.Request) {

	var (
		method = r.Method
		url    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	_, err := w.Write([]byte("pong"))
	if err != nil {
		app.logger.Error(err.Error(), "method", method, "url", url, "trace", trace)
		return
	}
}

// home renders the home template and responds with an HTTP 200 status. It does not take or process any additional data.
func (app *App) home(w http.ResponseWriter, r *http.Request) {

	events, err := app.eventModel.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.context.Events = events

	app.render(w, r, "home.tmpl", app.context, http.StatusOK)
}

// eventDetail retrieves the details of a specific event based on the ID from the URL, renders the detail template, and responds.
func (app *App) eventDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	event, err := app.eventModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Assign the retrieved event data to the application data structure.
	app.context.Event = event

	app.render(w, r, "events/detail.tmpl", app.context, http.StatusOK)
}

// eventCreate renders the "create event" template and responds with an HTTP 200 status. It does not process input data.
func (app *App) eventCreate(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "events/create.tmpl", app.context, http.StatusOK)
}

// eventCreatePost handles the POST request for creating an event, parses the form data, and validates the request.
func (app *App) eventCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest, err)
		return
	}

	var form EventForm

	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest, err)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field is required.")
	form.CheckField(validator.NotBlank(form.Location), "location", "This field is required.")
	form.CheckField(validator.NotBlank(form.Description), "description", "This field is required.")
	form.CheckField(validator.ValidDate(form.EventDate), "eventDate", "This field is required.")

	if !form.Valid() {
		app.Form = form
		app.render(w, r, "events/create.tmpl", app.context, http.StatusUnprocessableEntity)
		return
	}

	fmt.Println(form)
	http.Redirect(w, r, "/", http.StatusFound)
}
