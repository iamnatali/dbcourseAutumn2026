package queries

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct {
	pool *pgxpool.Pool
	logger *slog.Logger
}

func New(pgxPool *pgxpool.Pool, logger *slog.Logger) *Queries {
	return &Queries{pool: pgxPool, logger: logger}
}

