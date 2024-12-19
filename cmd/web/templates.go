package main

import (
	"bytes"
	"fmt"
	"github.com/madalinpopa/go-event-planner/ui"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// functions is a template.FuncMap containing custom template functions for use in HTML templates.
var functions = template.FuncMap{
	"humanDate": func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.UTC().Format("2006-01-02")
	},
}

// newTemplateCache initializes and returns a cache of precompiled templates, or an error if the operation fails.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Walk through all files in the embedded filesystem
	err := fs.WalkDir(ui.Files, "html/pages", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip if it's a directory or not a template file
		if d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// Generate cache key by removing "html/pages/" prefix and maintaining subdirectory structure
		name := strings.TrimPrefix(path, "html/pages/")

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			"html/partials/events/*.tmpl",
			path,
		}

		ts, err := template.New(filepath.Base(path)).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return fmt.Errorf("parsing template %s: %w", path, err)
		}

		cache[name] = ts
		fmt.Printf("Cached template: %s\n", name)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking pages directory: %w", err)
	}

	return cache, nil
}

// render writes a rendered template to the response writer with the given status code.
// It checks if the template exists, handles errors, and logs issues appropriately.
// If the template is successfully rendered, its output is written to the response.
func (app *App) render(w http.ResponseWriter, r *http.Request, name string, data interface{}, status int) {

	// Check if the template with the given name exists in the template cache.
	// If the template is not found, respond with a server error and stop further processing.
	t, ok := app.templates[name]

	if !ok {
		err := fmt.Errorf("template %s not found", name)
		app.serverError(w, r, err)
		return
	}

	// Create a new buffer to hold the rendered template.
	// The use of a buffer allows for efficient error handling and ensures that
	// the template output is only written to the response writer after successful processing.
	buf := new(bytes.Buffer)

	err := t.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
