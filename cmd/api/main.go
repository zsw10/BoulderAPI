package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"sync"

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
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.boulder.url, "boulderURL", "http://localhost:4001/directory", "boudler CA URL")
	flag.StringVar(&cfg.db.file, "db-file", "./testDB.db", "Sqlite DB file")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := sql.Open("sqlite3", cfg.db.file)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
