package redis

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		slog.Error("could not parse redis url", "error", err)
		return nil, errors.New("could not parse redis url")
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Error("could not ping redis", "error", err)
		return nil, errors.New("could not ping redis")
	}

	slog.Info("connected to redis")

	return client, nil
}
