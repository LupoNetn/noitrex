package events

import (
	"fmt"
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
