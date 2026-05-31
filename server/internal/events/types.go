package events

import (
	"fmt"
	"time"

	"github.com/luponetn/noitrex/internal/db"
)

//Errors both sentinel and custom error types

type UsageEventExists struct {
	IdempotencyKey string
	CustomerID     string
}

func (u *UsageEventExists) Error() string {
	return fmt.Sprintf("an event for this user with id %v already exists under this idempotency key %v", u.CustomerID, u.IdempotencyKey)
}

// Event Structures
type UsageEventRequest struct {
	CustomerID     string `json:"customer_id" binding:"required"`
	EventName      string `json:"event_name" binding:"required"`
	Quantity       int64  `json:"quantity" binding:"required"`
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
}

type UsageEventResponse struct {
	ID             string `json:"id"`
	CustomerID     string `json:"customer_id"`
	EventName      string `json:"event_name"`
	Quantity       int64  `json:"quantity"`
	IdempotencyKey string `json:"idempotency_key"`
	CreatedAt      string `json:"created_at"`
}

func toUsageEventResponse(e db.UsageEvent) UsageEventResponse {
	var idStr, customerIdStr string
	if e.ID.Valid {
		idStr = fmt.Sprintf("%x-%x-%x-%x-%x", e.ID.Bytes[0:4], e.ID.Bytes[4:6], e.ID.Bytes[6:8], e.ID.Bytes[8:10], e.ID.Bytes[10:16])
	}
	if e.CustomerID.Valid {
		customerIdStr = fmt.Sprintf("%x-%x-%x-%x-%x", e.CustomerID.Bytes[0:4], e.CustomerID.Bytes[4:6], e.CustomerID.Bytes[6:8], e.CustomerID.Bytes[8:10], e.CustomerID.Bytes[10:16])
	}

	return UsageEventResponse{
		ID:             idStr,
		CustomerID:     customerIdStr,
		EventName:      e.EventName,
		Quantity:       e.Quantity,
		IdempotencyKey: e.IdempotencyKey,
		CreatedAt:      e.CreatedAt.Time.Format(time.RFC3339),
	}
}
