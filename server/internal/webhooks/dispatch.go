package webhook

import (
	"context"
	"log/slog"
	"sync"

	"github.com/luponetn/nexusmq/pkg/broker"
	"github.com/luponetn/noitrex/internal/db"
)

type WebhookEngine struct {
	db     db.Querier
	broker broker.Broker
}

func NewWebhookEngine(q db.Querier, b broker.Broker) *WebhookEngine {
	return &WebhookEngine{
		db:     q,
		broker: b,
	}
}

func (w *WebhookEngine) Start(ctx context.Context) {
	const (
		bufferSize  = 512
		workerCount = 20
	)

	brokerchan := make(chan receiveResult, bufferSize)

	sub, err := w.broker.Subscribe("invoice.created")
	if err != nil {
		slog.Error("webhook engine failed to start", slog.String("error", err.Error()))
		return
	}

	// receiver goroutine
	go func() {
		for {
			msg, err := sub.Receive()
			select {
			case <-ctx.Done():
				close(brokerchan)
				return
			case brokerchan <- receiveResult{msg, err}:
			}
		}
	}()

	// fixed worker pool — no unbounded goroutine spawning
	var wg sync.WaitGroup
	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for res := range brokerchan {
				if res.err != nil {
					slog.Error("webhook engine receive error", slog.String("error", res.err.Error()))
					continue
				}
				w.HandleDispatchWebhook(res.msg)
			}
		}()
	}

	wg.Wait()
}

func (w *WebhookEngine) HandleDispatchWebhook(msg *broker.Message) {

}
