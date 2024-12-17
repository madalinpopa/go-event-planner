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

func (app *App) home(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("home"))
	if err != nil {
		panic(err)
	}
}
