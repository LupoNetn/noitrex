package metering

import (
	"context"
	"log/slog"
	"time"

	"github.com/luponetn/nexusmq/pkg/broker"
	"github.com/redis/go-redis/v9"
)

type MeteringEngine struct {
	broker      broker.Broker
	redisClient *redis.Client
}

func NewMeteringEngine(b broker.Broker, r *redis.Client) *MeteringEngine {
	return &MeteringEngine{
		broker:      b,
		redisClient: r,
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
				m.FlushMeterUsage(ctx)
				slog.Info("shuting down metering engine")
				return
			case res := <-resChan:
				if res.err != nil {
					slog.Error("failed to receive message", "error", res.err)
					continue
				}
				m.ProcessUsageEvent(res.msg)
			case <-ticker.C:
				m.FlushMeterUsage(ctx)
			}
		}
	}()
}

func (m *MeteringEngine) ProcessUsageEvent(msg *broker.Message) {

}

func (m *MeteringEngine) FlushMeterUsage(ctx context.Context) error {
	return nil
}
