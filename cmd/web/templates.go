package main

import (
	"github.com/madalinpopa/go-event-planner/ui"
	"html/template"
	"io/fs"
	"path/filepath"
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
	cache := make(map[string]*template.Template)

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/pages/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return nil, nil
}
