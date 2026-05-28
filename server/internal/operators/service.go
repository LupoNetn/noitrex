package operator

import (
	"github.com/luponetn/noitrex/internal/db"
)

type Service interface {
	//All functions here are placeholders and shoud be properly implemented soon
	GetOperatorByID(id string) (string, error)
	UpdateOperatorDetails(id string) (string, error)
}

type Svc struct {
	db db.Querier
}

func NewService(db db.Querier) Service {
	return &Svc{
		db: db,
	}
}

// implement servie interface
func (s *Svc) GetOperatorByID(id string) (string, error) {
	return "", nil
}

func (s *Svc) UpdateOperatorDetails(id string) (string, error) {
	return "", nil
}
