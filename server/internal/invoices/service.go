package invoices

import "github.com/luponetn/noitrex/internal/db"

type Service interface{}

type Svc struct {
	db db.Querier
}

func NewService(db db.Querier) Service {
	return &Svc{db: db}
}
