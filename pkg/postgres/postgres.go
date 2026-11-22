package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"hotel.booking.natali/internal/handlers"
)

func NewPool(db_config handlers.DbConfig, logger *slog.Logger) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(db_config.Dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(db_config.MaxOpenConns)
	// todo: parse
	// poolConfig.MaxConnIdleTime = db_config.MaxIdleTime
	poolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   &loggerAdapter{logger: logger},
		LogLevel: tracelog.LogLevelTrace,
	}

	//todo: сделать так, чтобы при отмене запроса, запрос к БД тоже отменялся, прерывался 
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

type loggerAdapter struct {
	logger *slog.Logger
}

// todo: start using context.Context to pass logs between goroutines
func (l *loggerAdapter) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	switch level {
	case tracelog.LogLevelTrace:
		l.logger.DebugContext(ctx, msg, slog.Any("data", data))
	case tracelog.LogLevelDebug:
		l.logger.DebugContext(ctx, msg, slog.Any("data", data))
	case tracelog.LogLevelInfo:
		l.logger.InfoContext(ctx, msg, slog.Any("data", data))
	case tracelog.LogLevelWarn:
		l.logger.WarnContext(ctx, msg, slog.Any("data", data))
	case tracelog.LogLevelError:
		l.logger.ErrorContext(ctx, msg, slog.Any("data", data))
	case tracelog.LogLevelNone:
		l.logger.ErrorContext(ctx, msg, slog.Any("data", data))
	}
}
