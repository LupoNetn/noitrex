package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/nexusmq/pkg/broker"
	"github.com/luponetn/noitrex/internal/db"
)

type WebhookEngine struct {
	db     db.Querier
	broker broker.Broker
	client *http.Client
}

func NewWebhookEngine(q db.Querier, b broker.Broker) *WebhookEngine {
	return &WebhookEngine{
		db:     q,
		broker: b,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
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
		defer close(brokerchan)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := sub.Receive()
				if err != nil {
					slog.Error("receive error", "error", err)
					return
				}
				brokerchan <- receiveResult{msg, err}
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
	var payload Invoice
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		slog.Error("webhook engine failed to unmarshal payload", slog.String("error", err.Error()))
		return
	}

	operatorId := payload.OperatorID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	webhookEndpoint, err := w.db.GetActiveWebhooksByOperatorID(ctx, operatorId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("webhook engine has no active webhook endpoint", slog.String("operator_id", operatorId.String()))
			return
		}
		slog.Error("webhook engine failed to get active webhook endpoint for operator",
			slog.String("operator_id", operatorId.String()),
			slog.String("error", err.Error()),
		)
		return
	}

	w.deliver(ctx, webhookEndpoint, payload)

}

func (w *WebhookEngine) deliver(ctx context.Context, endpoint db.WebhookEndpoint, payload Invoice) {
	//1.save a delivery attempt
	//2.send webhook
	//3. handle response if ok update delivery attempt row for succes
	///4. if not update delivery row with failure and check if error is retryable and try again
	//5. after 5 tries stop
	var lineItems InvoiceLineItems
	if err := json.Unmarshal(payload.LineItems, &lineItems); err != nil {
		slog.Error("webhook engine failed to unmarshal line items", slog.String("error", err.Error()))
		return
	}

	p, err := json.Marshal(payload)
	if err != nil {
		slog.Error("webhook engine failed to marshal payload", slog.String("error", err.Error()))
		return
	}

	delivery, err := w.db.CreateDeliveryAttempt(ctx, db.CreateDeliveryAttemptParams{
		EndpointID: endpoint.ID,
		EventType:  lineItems.Events[0],
		Payload:    p,
	})
	if err != nil {
		slog.Error("webhook engine failed to create delivery attempt", slog.String("error", err.Error()))
		return
	}

	w.send(ctx, endpoint, p, delivery.ID)

}

// webhook sender
func (w *WebhookEngine) send(
    ctx context.Context,
    endpoint db.WebhookEndpoint,
    payload []byte,
    deliveryID pgtype.UUID,
) error {
    const maxRetries = 6

    // sign once — signature doesn't change between retries
    mac := hmac.New(sha256.New, []byte(endpoint.SigningSecret))
    mac.Write(payload)
    signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))
    timestamp := fmt.Sprintf("%d", time.Now().Unix())

    var statusCode int
    var respBody []byte

    for attempt := 0; attempt < maxRetries; attempt++ {
        // bytes.NewReader supports re-reading on retry
        body := bytes.NewReader(payload)

        req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.Url, body)
        if err != nil {
            slog.Error("failed to create request", "error", err)
            return err
        }

        req.Header.Set("Content-Type",        "application/json")
        req.Header.Set("X-Noitrex-Signature", signature)
        req.Header.Set("X-Noitrex-Timestamp", timestamp)

        resp, err := w.client.Do(req)
        if err != nil {
            // network error — no response, retry
            slog.Error("webhook network error",
                "attempt", attempt+1,
                "error", err,
            )
            time.Sleep(retryDelay(attempt))
            continue
        }

        statusCode = resp.StatusCode
        respBody, _ = io.ReadAll(io.LimitReader(resp.Body, 1000))
        resp.Body.Close()

        // 2xx — success
        if statusCode >= 200 && statusCode < 300 {
            _, err := w.db.UpdateDeliveryResult(ctx, db.UpdateDeliveryResultParams{
                ID:           deliveryID,
                ResponseCode: pgtype.Int4{Int32: int32(statusCode), Valid: true},
                ResponseBody: pgtype.Text{String: string(respBody), Valid: true},
                Status:       "delivered",
            })
            if err != nil {
                slog.Error("failed to update delivery result", "error", err)
            }
            return nil
        }

        // 422 unprocessable — operator's endpoint rejected the payload
        // retrying won't help — stop immediately
        if statusCode == 422 {
            slog.Warn("webhook rejected by endpoint — not retrying",
                "status_code", statusCode,
                "endpoint", endpoint.Url,
            )
            break
        }

        // 5xx or other non-2xx — retry
        slog.Warn("webhook delivery failed",
            "attempt", attempt+1,
            "status_code", statusCode,
        )
        time.Sleep(retryDelay(attempt))
    }

    // exhausted all retries or got 422 — mark as failed
    _, err := w.db.UpdateDeliveryResult(ctx, db.UpdateDeliveryResultParams{
        ID:           deliveryID,
        ResponseCode: pgtype.Int4{Int32: int32(statusCode), Valid: true},
        ResponseBody: pgtype.Text{String: string(respBody), Valid: true},
        Status:       "failed",
    })
    if err != nil {
        slog.Error("failed to update delivery result as failed", "error", err)
        return err
    }

    slog.Error("webhook delivery permanently failed",
        "endpoint", endpoint.Url,
        "attempts", maxRetries,
    )
    return nil
}

// helper
func retryDelay(attemptNumber int) time.Duration {
	base := time.Duration(math.Pow(2, float64(attemptNumber))) * 30 * time.Second
	jitter := time.Duration(rand.Int63n(int64(base / 5)))
	return base + jitter
}
