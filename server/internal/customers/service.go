package customers

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/internal/db"
)

type Service interface {
	CreateCustomer(ctx context.Context, args db.CreateCustomerParams) (db.Customer, error)
	GetCustomerByID(ctx context.Context, id pgtype.UUID) (db.Customer, error)
	GetCustomerByExternalID(ctx context.Context, operatorID pgtype.UUID, externalID pgtype.UUID) (db.Customer, error)
	ListCustomers(ctx context.Context, operatorID pgtype.UUID) ([]db.Customer, error)
	GetCustomerByEmail(ctx context.Context, email string, operatorID pgtype.UUID) (db.Customer, error)
	GetCustomerInvoices(ctx context.Context, args db.GetCustomerInvoicesParams) ([]db.Invoice, error)
}

type Svc struct {
	db db.Querier
}

func NewService(db db.Querier) Service {
	return &Svc{
		db: db,
	}
}

func (s *Svc) CreateCustomer(ctx context.Context, args db.CreateCustomerParams) (db.Customer, error) {
	existingCustomer, err := s.db.GetCustomerByExternalID(ctx, db.GetCustomerByExternalIDParams{
		OperatorID: args.OperatorID,
		ExternalID: args.ExternalID,
	})
	if err == nil {
		return existingCustomer, &CustomerAlreadyExists{OperatorID: args.OperatorID.String(), ExternalID: args.ExternalID.String()}
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return existingCustomer, err
	}
	return s.db.CreateCustomer(ctx, args)

}

func (s *Svc) GetCustomerByID(ctx context.Context, id pgtype.UUID) (db.Customer, error) {
	return s.db.GetCustomerByID(ctx, id)
}

func (s *Svc) GetCustomerByExternalID(ctx context.Context, operatorID pgtype.UUID, externalID pgtype.UUID) (db.Customer, error) {
	return s.db.GetCustomerByExternalID(ctx, db.GetCustomerByExternalIDParams{
		OperatorID: operatorID,
		ExternalID: externalID,
	})
}

func (s *Svc) ListCustomers(ctx context.Context, operatorID pgtype.UUID) ([]db.Customer, error) {
	return s.db.ListCustomers(ctx, operatorID)
}

func (s *Svc) GetCustomerByEmail(ctx context.Context, email string, operatorID pgtype.UUID) (db.Customer, error) {
	return s.db.GetCustomerByEmail(ctx, db.GetCustomerByEmailParams{
		OperatorID: operatorID,
		Email:      email,
	})
}

func (s *Svc) GetCustomerInvoices(ctx context.Context, args db.GetCustomerInvoicesParams) ([]db.Invoice, error) {
	invoices, err := s.db.GetCustomerInvoices(ctx, db.GetCustomerInvoicesParams{
		OperatorID: args.OperatorID,
		CustomerID: args.CustomerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []db.Invoice{}, ErrNoCustomerInvoiceFound
		}
		return []db.Invoice{}, err
	}

	return invoices, nil
}
