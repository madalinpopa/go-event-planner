package main

import (
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	port string
)

type config struct {
	logger    *slog.Logger
	templates map[string]*template.Template
}

type App struct {
	config
}

func main() {

	flag.StringVar(&port, "port", "8080", "port to listen on")
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
