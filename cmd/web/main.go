package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	// port defines the default port the server listens on, typically overridden via a command-line flag.
	port string
)

// config is a struct that encapsulates application-wide dependencies, such as logging and template rendering.
type config struct {
	logger    *slog.Logger
	templates map[string]*template.Template
}

type data struct {
	Title       string
	CurrentYear int
}

// App is a struct that embeds configuration dependencies required across the application.
type App struct {
	config
	data
}

// main initializes the application, sets up dependencies, and starts the HTTP server listening on the specified port.
func main() {

	flag.StringVar(&port, "port", "4000", "port to listen on")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	templates, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := App{
		config: config{
			logger:    logger,
			templates: templates,
		},
		data: data{
			Title:       "Event Planner",
			CurrentYear: time.Now().Year(),
		},
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}

	logger.Info("Starting server", "port", port)

	if err := server.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
