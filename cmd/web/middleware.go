package main

import (
	"context"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

// addCommonHeaders is a middleware that adds common security-related HTTP headers to the response.
func (app *App) addCommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set Content-Security-Policy to restrict the sources of content such as scripts, styles, and images
		w.Header().Set("Content-Security-Policy", `
            default-src 'self';
            script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://*.iconify.design https://*.simplesvg.com https://*.unisvg.com;
            style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net;
            img-src 'self' https://cdn.jsdelivr.net https://*.iconify.design https://*.simplesvg.com https://*.unisvg.com;
            font-src 'self' https://cdn.jsdelivr.net;
            connect-src 'self' https://*.iconify.design https://*.simplesvg.com https://*.unisvg.com;
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

// csrfToken is a middleware that adds CSRF protection to HTTP handlers using the nosurf package.
func csrfToken(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

// loginRequired is a middleware that enforces authentication for protected routes.
// It redirects unauthenticated users to the login page and prevents caching of sensitive content.
func (app *App) loginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page and return
		// from the middleware chain so that no subsequent handlers in the chain are
		// executed.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		// And call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// redirectAuthenticatedUsers is a middleware that redirects authenticated users attempting
// to access the login or register pages to the events page.
func (app *App) redirectAuthenticatedUsers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authenticated using app's existing method
		if app.isAuthenticated(r) {
			// Check if the path is either the login or register page
			if r.URL.Path == "/login" || r.URL.Path == "/register" {
				// Redirect authenticated users to the events page
				http.Redirect(w, r, "/events", http.StatusSeeOther)
				return
			}
		}
		// If not authenticated or not targeting the restricted pages, continue
		next.ServeHTTP(w, r)
	})
}

// authenticate is a middleware that checks if a user is authenticated based on
// session data and updates the request data.
//
// If the user is not authenticated, the next handler in the chain is called without modifications
// to the request data.
//
// If the user is authenticated and exists in the
// database, the request data is updated to indicate the user's status.
func (app *App) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the authenticatedUserID value from the session using the GetInt()
		// method.
		//This will return the zero value for an int (0) if no
		// "authenticatedUserID" value is in the session -- in which case we call the
		// next handler in the chain as normal and return.
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Otherwise, we check to see if a user with that ID exists in our
		// database.
		exists, err := app.userModel.Exists(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		// If a matching user is found, we know that the request is coming from an
		// authenticated user who exists in our database.
		//We create a new copy of the
		// request (with an isAuthenticatedContextKey value of true in the request data)
		// and assign it to r.
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
