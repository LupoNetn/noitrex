package metering

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/nexusmq/pkg/broker"
	"github.com/luponetn/noitrex/internal/db"
	"github.com/redis/go-redis/v9"
)

type MeteringEngine struct {
	broker      broker.Broker
	redisClient *redis.Client
	db          db.Querier
}

func NewMeteringEngine(b broker.Broker, r *redis.Client, d db.Querier) *MeteringEngine {
	return &MeteringEngine{
		broker:      b,
		redisClient: r,
		db:          d,
	}
}

func (m *MeteringEngine) Start(ctx context.Context) {
	sub, err := m.broker.Subscribe("usage.ingested")
	if err != nil {
		slog.Error("failed to subscribe to usage.ingested", "error", err)
		return
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	resChan := make(chan receiveResult)

	go func() {
		for {
			msg, err := sub.Receive()
			resChan <- receiveResult{msg, err}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				// Use a fresh context for the final flush — the parent is already cancelled.
				flushCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				if err := m.FlushMeterUsage(flushCtx); err != nil {
					slog.Error("final flush failed on shutdown", "error", err)
				}
				cancel()
				slog.Info("shutting down metering engine")
				return
			case res := <-resChan:
				if res.err != nil {
					slog.Error("failed to receive message", "error", res.err)
					continue
				}
				slog.Info("it fucking works!!, nexusmq workssss!!!")
				m.ProcessUsageEvent(res.msg)
			case <-ticker.C:
				if err := m.FlushMeterUsage(ctx); err != nil {
					slog.Error("periodic flush failed", "error", err)
				}
			}
		}
	}()
}

func (m *MeteringEngine) ProcessUsageEvent(msg *broker.Message) {
	var event db.UsageEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		slog.Error("failed to unmarshal usage event", "error", err)
		return
	}

	period := event.CreatedAt.Time.Format("2006-01")
	key := fmt.Sprintf("usage:%s:%s:%s:%s",
		event.OperatorID.String(),
		event.CustomerID.String(),
		event.EventName,
		period,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: use pipeline for multiple increments
	if err := m.redisClient.IncrBy(ctx, key, event.Quantity).Err(); err != nil {
		slog.Error("failed to increment usage", "error", err)
		return
	}

	slog.Info("incremented counter", "key", key, "quantity", event.Quantity)
}

func (m *MeteringEngine) FlushMeterUsage(ctx context.Context) error {
	// Collect all usage:* keys via SCAN — non-blocking unlike KEYS.
	var (
		cursor uint64
		keys   []string
	)
	for {
		batch, nextCursor, err := m.redisClient.Scan(ctx, cursor, "usage:*", 100).Result()
		if err != nil {
			return fmt.Errorf("redis scan failed: %w", err)
		}
		keys = append(keys, batch...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	if len(keys) == 0 {
		return nil
	}

	slog.Info("flushing metering data", "keys", len(keys))

	for _, key := range keys {
		// Key format: usage:{operatorID}:{customerID}:{eventName}:{period}
		// SplitN to 5 keeps any ":" inside an event name intact.
		parts := strings.SplitN(key, ":", 5)
		if len(parts) != 5 {
			slog.Warn("skipping malformed metering key", "key", key)
			continue
		}

		operatorIDStr := parts[1]
		customerIDStr := parts[2]
		eventName := parts[3]
		period := parts[4] // "YYYY-MM"

		// GETDEL atomically reads and removes the key, preventing double-counting
		val, err := m.redisClient.GetDel(ctx, key).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			slog.Error("failed to getdel key", "key", key, "error", err)
			continue
		}

		totalUnits, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			slog.Error("failed to parse usage value", "key", key, "value", val, "error", err)
			continue
		}

		// Derive the calendar-month boundaries from the "YYYY-MM" period string.
		periodStart, err := time.Parse("2006-01", period)
		if err != nil {
			slog.Error("failed to parse period", "period", period, "error", err)
			continue
		}
		// period_end is exclusive — the first instant of the following month.
		periodEnd := periodStart.AddDate(0, 1, 0)

		var operatorID pgtype.UUID
		if err := operatorID.Scan(operatorIDStr); err != nil {
			slog.Error("failed to parse operator ID", "key", key, "error", err)
			continue
		}

		var customerID pgtype.UUID
		if err := customerID.Scan(customerIDStr); err != nil {
			slog.Error("failed to parse customer ID", "key", key, "error", err)
			continue
		}

		if err := m.db.UpsertUsageAggregate(ctx, db.UpsertUsageAggregateParams{
			OperatorID:  operatorID,
			CustomerID:  customerID,
			EventName:   eventName,
			PeriodStart: pgtype.Timestamptz{Time: periodStart, Valid: true},
			PeriodEnd:   pgtype.Timestamptz{Time: periodEnd, Valid: true},
			TotalUnits:  totalUnits,
		}); err != nil {
			slog.Error("failed to upsert usage aggregate", "key", key, "error", err)
			continue
		}

		slog.Info("flushed usage aggregate", "key", key, "total_units", totalUnits)
	}

	return nil
}
