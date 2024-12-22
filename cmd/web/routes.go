package main

import (
	"github.com/justinas/alice"
	"github.com/madalinpopa/go-event-planner/ui"
	"net/http"
)

// routes initializes and returns the HTTP router with predefined application routes and their handlers.
func (app *App) routes() http.Handler {

	// Create a new HTTP request multiplexer to handle application routing.
	mux := http.NewServeMux()

	// Serve static files from the embedded filesystem under the /static/ path.
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// Use the nosurf middleware on all 'csrfProtect' routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave, csrfToken, app.authenticate)

	protected := dynamic.Append(app.loginRequired)

	// Public routes
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /ping", dynamic.ThenFunc(app.ping))
	mux.Handle("GET /events/detail/{id}", dynamic.ThenFunc(app.eventView))
	mux.Handle("GET /events", dynamic.ThenFunc(app.eventList))

	// Protected routes
	mux.Handle("GET /events/create", protected.ThenFunc(app.eventCreate))
	mux.Handle("POST /events/create", protected.ThenFunc(app.eventCreatePost))
	mux.Handle("POST /events/{id}/delete", protected.ThenFunc(app.eventDelete))

	// User registration and authentication routes
	mux.Handle("GET /login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("GET /register", dynamic.ThenFunc(app.userRegister))
	mux.Handle("POST /register", dynamic.ThenFunc(app.userRegisterPost))
	mux.Handle("POST /logout", dynamic.ThenFunc(app.userLogoutPost))

	// Initialize middleware chain with panic recovery, request logging, and common headers.
	standardMiddleware := alice.New(app.addPanicRecover, app.addRequestLogger, app.addCommonHeaders)

	return standardMiddleware.Then(mux)
}
