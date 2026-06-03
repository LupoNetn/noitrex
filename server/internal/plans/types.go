package plans

import (
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/internal/db"
)

//Errors

var (
	ErrPlanAlreadyExists = errors.New("plan already exists")
	ErrPlanNotFound      = errors.New("plan not found")
)

// Types
type CreatePlanRequest struct {
	Name          string               `json:"name" binding:"required"`
	PricingModel  db.PricingType       `json:"pricing_model" binding:"required"`
	UnitPriceCent pgtype.Int8          `json:"unit_price_cent" binding:"required"`
	Tiers         []byte               `json:"tiers"`
	BillingPeriod db.BillingPeriodType `json:"billing_period" binding:"required"`
}

type UpdatePlanRequest struct {
	ID            pgtype.UUID          `json:"id" binding:"required"`
	Name          string               `json:"name"`
	PricingModel  db.PricingType       `json:"pricing_model"`
	UnitPriceCent pgtype.Int8          `json:"unit_price_cent"`
	Tiers         []byte               `json:"tiers"`
	BillingPeriod db.BillingPeriodType `json:"billing_period"`
}
