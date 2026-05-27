package db

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBPool(dbUrl string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		slog.Error("failed to parse db url", "error", err.Error())
		return nil, err
	}

	config.MaxConns = 20
	config.MinConns = 5
	config.HealthCheckPeriod = 10 * time.Second
	config.ConnConfig.RuntimeParams = map[string]string{
		"application_name": "noitrex",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("failed to parse db url", "error", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		slog.Error("db ping failed", "error", err)
		return nil, err
	}

	slog.Info("db connection established successfully")

	return pool, nil
}
