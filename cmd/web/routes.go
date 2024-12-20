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

	// Additional routes can be added here as needed.
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /ping", app.ping)
	mux.HandleFunc("GET /events/detail/{id}", app.eventDetail)
	mux.HandleFunc("GET /events/create", app.eventCreate)
	mux.HandleFunc("POST /events/create", app.eventCreatePost)

	// Initialize middleware chain with panic recovery, request logging, and common headers.
	standardMiddleware := alice.New(app.addPanicRecover, app.addRequestLogger, app.addCommonHeaders)

	return standardMiddleware.Then(mux)
}
