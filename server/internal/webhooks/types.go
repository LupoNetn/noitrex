package webhook

import (
	"errors"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/luponetn/nexusmq/pkg/broker"
	"github.com/luponetn/noitrex/internal/db"
)

// Errors
var (
	ErrWebhookEndpointNotFound = errors.New("webhook endpoint not found")
)

// Types
type receiveResult struct {
	msg *broker.Message
	err error
}

type Invoice struct {
	ID          pgtype.UUID
	OperatorID  pgtype.UUID
	CustomerID  pgtype.UUID
	AmountCents int64
	Status      db.InvoiceStatus
	PeriodStart pgtype.Timestamptz
	PeriodEnd   pgtype.Timestamptz
	LineItems   []byte
	CreatedAt   pgtype.Timestamptz
}

type InvoiceLineItems struct {
	Events      []string `json:"events"`
	TotalUnits  int64    `json:"total_units"`
	TotalAmount int64    `json:"total_amount"`
	PeriodStart string   `json:"period_start"`
	PeriodEnd   string   `json:"period_end"`
	Plan        string   `json:"plan"`
}

type PaginatedRequest struct {
	Limit  int32 `form:"limit,default=10"`
	Offset int32 `form:"offset,default=0"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	TotalCount int64       `json:"total_count"`
	Limit      int32       `json:"limit"`
	Offset     int32       `json:"offset"`
}

type UpdateEndpointUrlRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type WebhookDeliveryStatsResponse struct {
	TotalDeliveries      int64 `json:"total_deliveries"`
	SuccessfulDeliveries int64 `json:"successful_deliveries"`
	FailedDeliveries     int64 `json:"failed_deliveries"`
}
