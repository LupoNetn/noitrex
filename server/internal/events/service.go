package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/luponetn/nexusmq/pkg/broker"
	"github.com/luponetn/noitrex/internal/db"
)

type Service interface {
	CreateUsageEvent(ctx context.Context, args db.CreateUsageEventParams) (db.UsageEvent, error)
}

type Svc struct {
	db     db.Querier
	broker broker.Broker
}

func NewService(db db.Querier, b broker.Broker) Service {
	return &Svc{
		db:     db,
		broker: b,
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

	payload, err := json.Marshal(newEvent)
	if err != nil {
		slog.Error("failed to marshal usage event for broker", "error", err)
		return newEvent, nil
	}

	if err := s.broker.Publish("usage.ingested", &broker.Message{
		Payload:   payload,
		Timestamp: time.Now(),
	}); err != nil {
		slog.Error("failed to publish usage event to broker", "error", err)
	}

	return newEvent, nil

}
