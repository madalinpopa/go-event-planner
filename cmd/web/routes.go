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
	csrfProtect := alice.New(app.sessionManager.LoadAndSave, csrfToken)

	// Additional routes can be added here as needed.
	mux.Handle("GET /{$}", csrfProtect.ThenFunc(app.home))
	mux.Handle("GET /ping", csrfProtect.ThenFunc(app.ping))
	mux.Handle("GET /events/detail/{id}", csrfProtect.ThenFunc(app.eventDetail))
	mux.Handle("GET /events/create", csrfProtect.ThenFunc(app.eventCreate))
	mux.Handle("POST /events/create", csrfProtect.ThenFunc(app.eventCreatePost))
	mux.Handle("POST /events/{id}/delete", csrfProtect.ThenFunc(app.eventDelete))

	// Initialize middleware chain with panic recovery, request logging, and common headers.
	standardMiddleware := alice.New(app.addPanicRecover, app.addRequestLogger, app.addCommonHeaders)

	return standardMiddleware.Then(mux)
}
