package main

import (
	"database/sql"
	"flag"
	"github.com/madalinpopa/go-event-planner/internal/models"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	// port defines the default port the server listens on, typically overridden via a command-line flag.
	port string
)

// config is a struct that encapsulates application-wide dependencies, such as logging and template rendering.
type config struct {
	logger    *slog.Logger
	templates map[string]*template.Template
	db        *sql.DB
}

type context struct {
	Title       string
	CurrentYear int
	Events      []models.Event
	Event       models.Event
}

// App is a struct that embeds configuration dependencies required across the application.
type App struct {
	eventModel *models.EventModel
	config
	context
}

// main initializes the application, sets up dependencies, and starts the HTTP server listening on the specified port.
func main() {

	flag.StringVar(&port, "port", "4000", "port to listen on")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB("database/events.db")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	templates, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := App{
		eventModel: &models.EventModel{DB: db},
		config: config{
			logger:    logger,
			templates: templates,
			db:        db,
		},
		context: context{
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
