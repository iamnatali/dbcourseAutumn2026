package queries

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"hotel.booking.natali/internal/domain"
)

const searchMoviesQuery = `SELECT * FROM movies;`
const getMovieQuery = `SELECT * FROM movies WHERE id = $1;`
const insertMoviesQuery = `INSERT INTO movies (title, release_date, runtime, genres, language) 
VALUES ($1, $2, $3, $4, $5) RETURNING id;`

const updateMovie = `UPDATE movies
SET (title, release_date, runtime, genres, language) = ($1, $2, $3, $4, $5)
WHERE id = $6;`
const deleteMovie = `DELETE FROM movies WHERE id = $1 RETURNING *;`

type DbMovie struct {
	ID          uuid.UUID       `db:"id"`
	Title       string          `db:"title"`
	ReleaseDate pgtype.Date     `db:"release_date"`
	Runtime     pgtype.Interval `db:"runtime"`
	Genres      []string        `db:"genres"`
	Language    string          `db:"language"`
}

func (q *Queries) DeleteMovie(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	rows, err := q.pool.Query(ctx, deleteMovie, id)
	dbMovie, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DbMovie])
	if err != nil {
		q.logger.Error("Database query error happenned during movie deletion",
			slog.Any("error", err))
		return nil, err
	}

	deleted_movie := domain.Movie{
		ID:          dbMovie.ID,
		Title:       dbMovie.Title,
		ReleaseDate: domain.Date(dbMovie.ReleaseDate.Time),
		Runtime:     domain.Runtime((time.Duration(dbMovie.Runtime.Microseconds) * time.Microsecond).Minutes()),
		Genres:      dbMovie.Genres,
		Language:    dbMovie.Language,
	}
	q.logger.Info("Movie was successfully deleted", slog.Any("deleted_movie", deleted_movie))
	return &deleted_movie, nil
}

// nice to have: add genres
// todo: generalize error handling
func (q *Queries) UpdateMovie(ctx context.Context, movie *domain.Movie) error {
	pgDate := pgtype.Date{
		Time:  time.Time(movie.ReleaseDate),
		Valid: true,
	}
	args := []interface{}{movie.Title, pgDate, fmt.Sprintf("%d minutes", movie.Runtime),
		pgtype.FlatArray[string](movie.Genres), movie.Language, movie.ID}

	commandTag, err := q.pool.Exec(ctx, updateMovie, args...)
	if err != nil {
		q.logger.Error("Database query error happenned during movie update",
			slog.Any("error", err), slog.Any("update_movie", movie))
		return err
	}

	if commandTag.RowsAffected() == 0 {
		q.logger.Error("Database query error happenned during movie update", slog.Any("error", err), slog.Any("movie", movie))
	} else {
		q.logger.Info("Movie was successfully updated", slog.Any("movie", movie))
	}

	return nil
}

func (q *Queries) InsertMovie(ctx context.Context, movie *domain.Movie) error {
	pgDate := pgtype.Date{
		Time:  time.Time(movie.ReleaseDate),
		Valid: true,
	}
	args := []interface{}{movie.Title, pgDate, fmt.Sprintf("%d minutes", movie.Runtime),
		pgtype.FlatArray[string](movie.Genres), movie.Language}

	row := q.pool.QueryRow(ctx, insertMoviesQuery, args...)
	err := row.Scan(&movie.ID)

	if err != nil {
		q.logger.Error("Database query error happenned during movie insert", slog.Any("error", err))
		return err
	}

	q.logger.Info("Movie was successfully inserted")
	return nil
}

const selectError = "Database query error happenned during movies select"

func (q *Queries) GetMovie(ctx context.Context, id uuid.UUID) (domain.Movie, error) {
	rows, err := q.pool.Query(ctx, getMovieQuery, id)
	var movie domain.Movie
	if err != nil {
		q.logger.Error("Database query error happenned during movie get", slog.Any("error", err))
		return movie, err
	}

	dbMovie, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[DbMovie])

	if err != nil {
		q.logger.Error("Database query error happenned during movie get", slog.Any("error", err))
		return movie, err
	}

	return domain.Movie{
		ID:          dbMovie.ID,
		Title:       dbMovie.Title,
		ReleaseDate: domain.Date(dbMovie.ReleaseDate.Time),
		Runtime:     domain.Runtime((time.Duration(dbMovie.Runtime.Microseconds) * time.Microsecond).Minutes()),
		Genres:      dbMovie.Genres,
		Language:    dbMovie.Language,
	}, nil
}

func (q *Queries) SearchMovies(ctx context.Context) ([]domain.Movie, error) {
	rows, err := q.pool.Query(ctx, searchMoviesQuery)
	if err != nil {
		q.logger.Error(selectError, slog.Any("error", err))
		return nil, err
	}

	defer rows.Close()

	var candles []domain.Movie
	for rows.Next() {
		var c domain.Movie
		var runtimeInt int32
		err = rows.Scan(&c.ID, &c.Title, &c.ReleaseDate, &runtimeInt, &c.Genres, &c.Language)
		c.Runtime = domain.Runtime(runtimeInt)
		if err != nil {
			return nil, err
		}
		candles = append(candles, c)
	}
	if err = rows.Err(); err != nil {
		q.logger.Error(selectError, slog.Any("error", err))
		return nil, err
	}

	return candles, nil
}
