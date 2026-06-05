package invoices

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/internal/db"
)

type Service interface {
	GetInvoiceByID(ctx context.Context, opID pgtype.UUID, invoiceID pgtype.UUID) (db.Invoice, error)
	ListInvoicesByCustomer(ctx context.Context, opID pgtype.UUID, customerID pgtype.UUID) ([]db.Invoice, error)
	ListOperatorInvoices(ctx context.Context, opID pgtype.UUID, limit, offset int32) ([]db.Invoice, int, error)
	UpdateInvoiceStatus(ctx context.Context, opID pgtype.UUID, invoiceID pgtype.UUID, status db.InvoiceStatus) (db.Invoice, error)
}

type Svc struct {
	db db.Querier
}

func NewService(db db.Querier) Service {
	return &Svc{db: db}
}

func (s *Svc) GetInvoiceByID(ctx context.Context, opID pgtype.UUID, invoiceID pgtype.UUID) (db.Invoice, error) {
	inv, err := s.db.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return db.Invoice{}, err
	}
	if inv.OperatorID != opID {
		return db.Invoice{}, ErrUnauthorized
	}
	return inv, nil
}

func (s *Svc) ListInvoicesByCustomer(ctx context.Context, opID pgtype.UUID, customerID pgtype.UUID) ([]db.Invoice, error) {
	invoices, err := s.db.ListInvoicesByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}
	var filtered []db.Invoice
	for _, inv := range invoices {
		if inv.OperatorID == opID {
			filtered = append(filtered, inv)
		}
	}
	return filtered, nil
}

func (s *Svc) UpdateInvoiceStatus(ctx context.Context, opID pgtype.UUID, invoiceID pgtype.UUID, status db.InvoiceStatus) (db.Invoice, error) {
	inv, err := s.db.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return db.Invoice{}, err
	}
	if inv.OperatorID != opID {
		return db.Invoice{}, ErrUnauthorized
	}

	return s.db.UpdateInvoiceStatus(ctx, db.UpdateInvoiceStatusParams{
		ID:     invoiceID,
		Status: status,
	})
}

func (s *Svc) ListOperatorInvoices(ctx context.Context, opID pgtype.UUID, limit, offset int32) ([]db.Invoice, int, error) {
	invList, err := s.db.ListInvoicesByOperatorPaginated(ctx, db.ListInvoicesByOperatorPaginatedParams{
		OperatorID: opID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, 0, err
	}
	// Total count: we can compute based on offset + returned rows for now.
	// A separate COUNT query can be added later if needed.
	total := int(offset) + len(invList)
	return invList, total, nil
}
