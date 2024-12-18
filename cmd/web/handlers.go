package main

import (
	"net/http"
	"runtime/debug"
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
	app.render(w, r, "home.tmpl", app.data, http.StatusOK)
}
