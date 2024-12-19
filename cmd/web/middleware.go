package main

import (
	"fmt"
	"net/http"
)

// addCommonHeaders is a middleware that adds common security-related HTTP headers to the response.
func (app *App) addCommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set Content-Security-Policy to restrict the sources of content such as scripts, styles, and images
		w.Header().Set("Content-Security-Policy", `
			default-src 'self';
			script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net;
			style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net;
			img-src 'self' https://cdn.jsdelivr.net;
			font-src 'self' https://cdn.jsdelivr.net;
		`)
		// Set Referrer-Policy to control the amount of referrer information sent with requests
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		// Set X-Content-Type-Options to prevent browsers from interpreting files as a different MIME type
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Set X-Frame-Options to prevent clickjacking attacks by disallowing the page from being framed
		w.Header().Set("X-Frame-Options", "deny")

		// Set X-XSS-Protection to disable the browser's XSS protection, preventing unintended behavior
		w.Header().Set("X-XSS-Protection", "0")

		// Set Server header to specify the server software being used (can also be customized or omitted)
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

// addRequestLogger logs details of incoming HTTP requests, such as IP, protocol, method, and URL, before passing to the next handler.
func (app *App) addRequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			url    = r.URL.RequestURI()
		)

		app.logger.Info("request", "ip", ip, "proto", proto, "method", method, "url", url)

		next.ServeHTTP(w, r)
	})
}

// addPanicRecover adds middleware to recover from panics during request handling and responds with a server error.
func (app *App) addPanicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}