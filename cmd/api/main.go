package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zsw10/BoulderAPI/internal/data"
)

type config struct {
	port    int
	boulder struct {
		url string
	}
	db struct {
		file string
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
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.boulder.url, "boulderURL", "http://localhost:4001/directory", "boudler CA URL")
	flag.StringVar(&cfg.db.file, "db-file", "./testDB.db", "Sqlite DB file")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := sql.Open("sqlite3", cfg.db.file)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
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
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
