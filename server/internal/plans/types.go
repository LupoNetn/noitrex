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
type PlanTiers struct {
	From  int64  `json:"from"`
	To    *int64 `json:"to"`
	Price int64  `json:"price"`
}

type CreatePlanRequest struct {
	Name          string               `json:"name" binding:"required"`
	PricingModel  db.PricingType       `json:"pricing_model" binding:"required"`
	UnitPriceCent pgtype.Int8          `json:"unit_price_cent" binding:"required"`
	Tiers         []PlanTiers          `json:"tiers"`
	BillingPeriod db.BillingPeriodType `json:"billing_period" binding:"required"`
}

type UpdatePlanRequest struct {
	Name          *string               `json:"name"`
	PricingModel  *db.PricingType       `json:"pricing_model"`
	UnitPriceCent *pgtype.Int8          `json:"unit_price_cent"`
	Tiers         *[]PlanTiers          `json:"tiers"`
	BillingPeriod *db.BillingPeriodType `json:"billing_period"`
}
