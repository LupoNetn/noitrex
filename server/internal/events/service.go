package events

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/luponetn/noitrex/internal/db"
)

type Service interface {
	CreateUsageEvent(ctx context.Context, args db.CreateUsageEventParams) (db.UsageEvent, error)
}

type Svc struct {
	db db.Querier
}

func NewService(db db.Querier) Service {
	return &Svc{
		db: db,
	}
}

func (s *Svc) CreateUsageEvent(ctx context.Context, args db.CreateUsageEventParams) (db.UsageEvent, error) {
	// check idempotency key
	event, err := s.db.GetUsageEventByIdempotencyKey(ctx, args.IdempotencyKey)
	if err == nil {
		// row found — duplicate
		slog.Info("duplicate usage event", "idempotency_key", args.IdempotencyKey)
		return event, &UsageEventExists{
			IdempotencyKey: args.IdempotencyKey,
			CustomerID:     args.CustomerID.String(),
		}
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("unexpected error checking idempotency key", "error", err)
		return db.UsageEvent{}, fmt.Errorf("check idempotency key: %w", err)
	}

	// safe to create — no duplicate found
	newEvent, err := s.db.CreateUsageEvent(ctx, args)
	if err != nil {
		slog.Error("failed to create usage event", "error", err)
		return db.UsageEvent{}, fmt.Errorf("create usage event: %w", err)
	}

	return newEvent, nil
}
