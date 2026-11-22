package handlers

import (
	"hotel.booking.natali/internal/repository"
	"log/slog"
)

const version = "1.0.0"

type DbConfig struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleTime  string
}

type Config struct {
	Port int
	Env  string
	Db   DbConfig
}

type Application struct {
	Config     Config
	Logger     *slog.Logger
	Repository *repository.Repository
}
