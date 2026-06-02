package customers

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/luponetn/noitrex/internal/db"
)

type Service interface {
	CreateCustomer(ctx context.Context, args db.CreateCustomerParams) (db.Customer, error)
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
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return existingCustomer, err
	}
	return s.db.CreateCustomer(ctx, args)

}
