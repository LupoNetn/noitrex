package auth

import (
	"github.com/luponetn/noitrex/internal/db"
)

type Service interface {
	//All functions here are placeholders and shoud be properly implemented soon
	Login(email string, password string) (string, error)
	Register(email string, password string) (string, error)
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
func (s *Svc) Login(email string, password string) (string, error) {
	return "", nil
}

func (s *Svc) Register(email string, password string) (string, error) {
	return "", nil
}
