package main

import (
	"net/http"
	"runtime/debug"
)

// serverError logs an internal server error and sends a 500 status response with a generic error message to the client.
func (app *App) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		url    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "url", url, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError logs client-side errors and sends the corresponding HTTP status code and message to the client.
func (app *App) clientError(w http.ResponseWriter, r *http.Request, status int, err error) {
	var (
		method = r.Method
		url    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method", method, "url", url, "status", status)
	http.Error(w, http.StatusText(status), status)
}

func (app *App) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
