package invoices

import (
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/internal/db"
)

// Errors
var ErrUnauthorized = errors.New("unauthorized access to invoice")

type UpdateInvoiceRequest struct {
	Status db.InvoiceStatus `json:"status" binding:"required"`
}

type PaginatedInvoicesResponse struct {
	Invoices []*InvoiceResponse `json:"invoices"`
	Page     int                `json:"page"`
	Limit    int                `json:"limit"`
	Total    int                `json:"total"`
}

type InvoiceResponse struct {
	ID          pgtype.UUID        `json:"id"`
	OperatorID  pgtype.UUID        `json:"operator_id"`
	CustomerID  pgtype.UUID        `json:"customer_id"`
	AmountCents int64              `json:"amount_cents"`
	Status      db.InvoiceStatus   `json:"status"`
	PeriodStart pgtype.Timestamptz `json:"period_start"`
	PeriodEnd   pgtype.Timestamptz `json:"period_end"`
	LineItems   map[string]any     `json:"line_items"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
}

func ToInvoiceResponse(inv db.Invoice) *InvoiceResponse {
	var lineItems map[string]any
	if len(inv.LineItems) > 0 {
		_ = json.Unmarshal(inv.LineItems, &lineItems)
	}

	return &InvoiceResponse{
		ID:          inv.ID,
		OperatorID:  inv.OperatorID,
		CustomerID:  inv.CustomerID,
		AmountCents: inv.AmountCents,
		Status:      inv.Status,
		PeriodStart: inv.PeriodStart,
		PeriodEnd:   inv.PeriodEnd,
		LineItems:   lineItems,
		CreatedAt:   inv.CreatedAt,
	}
}
