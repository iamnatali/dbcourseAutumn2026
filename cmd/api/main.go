package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"hotel.booking.natali/internal/handlers"
	"hotel.booking.natali/internal/repository"
	"hotel.booking.natali/pkg/postgres"
)

func main() {
	var cfg handlers.Config
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment(development|staging|production)")
	//todo: take from os.Getenv("DB_DSN")
	flag.StringVar(&cfg.Db.Dsn, "db-dsn",
		"postgres://postgres:12345@localhost/cinemanetwork", "PostgreSQL DSN")
	flag.IntVar(&cfg.Db.MaxOpenConns,
		"db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.StringVar(&cfg.Db.MaxIdleTime,
		"db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	pg_pool, pool_creation_error := postgres.NewPool(cfg.Db, logger)
	if pool_creation_error != nil {
		logger.Error("Postgres pool creation error",
			slog.Any("error", pool_creation_error))
	}

	defer pg_pool.Close()

	repo := repository.NewRepository(pg_pool, logger)

	app := &handlers.Application{
		Config:     cfg,
		Logger:     logger,
		Repository: &repo,
	}

	r := app.Routes()
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           r,
		IdleTimeout:       time.Minute,
		ReadTimeout:       10 * time.Second, //todo: move read timeout behavior to handlers
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      30 * time.Second,
	}

	logger.Info("Starting %s http server on %s", cfg.Env, server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		logger.Error("Http server error", slog.Any("error", err))
	}
}
