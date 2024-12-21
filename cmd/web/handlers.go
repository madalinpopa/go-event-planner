package main

import (
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
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
	app.context.CSRFToken = nosurf.Token(r)
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
	app.Form = EventForm{}
	app.render(w, r, "events/create.tmpl", app.context, http.StatusOK)
}

// eventDelete handles the deletion of an event record based on the ID extracted from the URL path.
// It returns a 404 Not Found error if the ID is invalid or the event does not exist.
// Logs internal errors and redirects to the home page upon successful deletion.
func (app *App) eventDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = app.eventModel.Delete(id)
	if err != nil {
		app.logger.Error(err.Error())
		app.serverError(w, r, err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	form.CheckField(validator.ValidDate(form.EventDate), "eventDate", "This field is required.")

	if !form.Valid() {
		fmt.Println(form)
		app.Form = form
		app.render(w, r, "events/create.tmpl", app.context, http.StatusUnprocessableEntity)
		return
	}

	id, err := app.eventModel.Create(form.Title, form.Description, form.EventDate, form.Location)
	fmt.Println(id, err)

	http.Redirect(w, r, "/", http.StatusFound)
}

// userRegister serves the user registration page by rendering the "register.tmpl"
// template with the application context.
func (app *App) userRegister(w http.ResponseWriter, r *http.Request) {
	app.context.Form = UserRegisterForm{}
	app.context.CSRFToken = nosurf.Token(r)
	app.render(w, r, "register.tmpl", app.context, http.StatusOK)
}

// userRegisterPost handles HTTP POST requests for user registration
// and renders the registration template with the given context.
func (app *App) userRegisterPost(w http.ResponseWriter, r *http.Request) {
	app.context.CSRFToken = nosurf.Token(r)

	var form UserRegisterForm

	err := app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.logger.Error(err.Error())
		app.clientError(w, r, http.StatusBadRequest, err)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field is required.")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field is required.")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "The email address is not valid.")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field is required.")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters.")

	if !form.Valid() {
		app.Form = form
		fmt.Println(form.FieldErrors)
		app.render(w, r, "register.tmpl", app.context, http.StatusUnprocessableEntity)
		return
	}

	err = app.userModel.Create(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "This email address is already registered.")
			app.Form = form
			app.render(w, r, "register.tmpl", app.context, http.StatusUnprocessableEntity)
		} else {
			app.logger.Error(err.Error())
			app.serverError(w, r, err)
		}
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// userLogin handles the user login page rendering by serving the login template
// with the appropriate context and status.
func (app *App) userLogin(w http.ResponseWriter, r *http.Request) {
	app.context.Form = UserLoginForm{}
	app.context.CSRFToken = nosurf.Token(r)
	app.render(w, r, "login.tmpl", app.context, http.StatusOK)
}

// userLoginPost handles POST requests for user login, rendering the login page
// with the provided context and HTTP status OK.
func (app *App) userLoginPost(w http.ResponseWriter, r *http.Request) {
	app.context.CSRFToken = nosurf.Token(r)

	form := UserLoginForm{}
	err := app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.logger.Error(err.Error())
		app.clientError(w, r, http.StatusBadRequest, err)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field is required.")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "The email address is not valid.")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field is required.")

	if !form.Valid() {
		app.Form = form
		app.render(w, r, "login.tmpl", app.context, http.StatusUnprocessableEntity)
		return
	}

	id, err := app.userModel.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Invalid email or password.")
			app.Form = form
			app.render(w, r, "login.tmpl", app.context, http.StatusUnprocessableEntity)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *App) userLogoutPost(w http.ResponseWriter, r *http.Request) {

	// Renew session token
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Remove authenticated user id
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
