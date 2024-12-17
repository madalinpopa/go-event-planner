package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	port string
)

type config struct {
	logger *slog.Logger
}

type App struct {
	config
}

func main() {

	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := App{
		config: config{
			logger: logger,
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
