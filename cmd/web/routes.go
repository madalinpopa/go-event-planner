package main

import "net/http"

// routes initializes and returns the HTTP router with predefined application routes and their handlers.
func (app *App) routes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", app.ping)
	mux.HandleFunc("GET /{$}", app.home)

	return mux
}
