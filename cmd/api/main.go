package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type config struct {
	port    int
	boulder struct {
		url string
	}
	db struct {
		file         string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}

	limiter struct {
		rps     int
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.Parse()
	cfg.boulder.url = "http://localhost:4001/directory"

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr)
	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
