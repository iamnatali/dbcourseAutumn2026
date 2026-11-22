package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"hotel.booking.natali/internal/domain"
	"hotel.booking.natali/internal/repository/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type repo struct {
	*queries.Queries
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewRepository(pgxPool *pgxpool.Pool, logger *slog.Logger) Repository {
	return &repo{
		Queries: queries.New(pgxPool, logger),
		pool:    pgxPool,
		logger:  logger,
	}
}

type Repository interface {
	DeleteMovie(ctx context.Context, id uuid.UUID) (*domain.Movie, error)
	UpdateMovie(ctx context.Context, movie *domain.Movie) error
	GetMovie(ctx context.Context, id uuid.UUID) (domain.Movie, error)
	SearchMovies(context.Context) ([]domain.Movie, error)
	InsertMovie(ctx context.Context, movie *domain.Movie) error
}
