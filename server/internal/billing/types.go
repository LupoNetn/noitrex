package billing

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/internal/db"
)

type UsageAggregate struct {
	EventName  string
	TotalUnits int64
}

type CustomerForBilling struct {
	CustomerID        pgtype.UUID          `json:"customer_id"`
	OperatorID        pgtype.UUID          `json:"operator_id"`
	CustomerName      string               `json:"customer_name"`
	PeriodStart       pgtype.Timestamptz   `json:"period_start"`
	PlanName          string               `json:"plan_name"`
	PlanPricingModel  db.PricingType       `json:"plan_pricing_model"`
	PlanTiers         []byte               `json:"plan_tiers"`
	PlanUnitPrice     pgtype.Int8          `json:"plan_unit_price"`
	PlanBillingPeriod db.BillingPeriodType `json:"plan_billing_period"`
	Aggregates        []UsageAggregate     `json:"aggregates"`
}
