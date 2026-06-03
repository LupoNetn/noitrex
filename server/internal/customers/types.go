package customers

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/internal/db"
)

// Errors Both Sentinel and Custom error types
var ErrCustomerNotFound = errors.New("customer not found")

type CustomerAlreadyExists struct {
	OperatorID string
	ExternalID string
}

func (c *CustomerAlreadyExists) Error() string {
	return fmt.Sprintf("could not register the customer %s for the operator %s as they already exist", c.OperatorID, c.ExternalID)
}

// Types
type CreateCustomerRequest struct {
	ExternalID string `json:"external_id" validate:"required"`
	Email      string `json:"email" validate:"required"`
	Name       string `json:"name" validate:"required"`
	PlanName   string `json:"plan_name" validate:"required"`
}

type CustomerResponse struct {
	ID         pgtype.UUID        `json:"id"`
	ExternalID pgtype.UUID        `json:"external_id"`
	Email      string             `json:"email"`
	Name       string             `json:"name"`
	PlanName   string             `json:"plan_name"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
}

func ToCustomerResponse(customer db.Customer) *CustomerResponse {
	return &CustomerResponse{
		ID:         customer.ID,
		ExternalID: customer.ExternalID,
		Email:      customer.Email,
		Name:       customer.Name,
		PlanName:   customer.PlanName,
		CreatedAt:  customer.CreatedAt,
	}
}
