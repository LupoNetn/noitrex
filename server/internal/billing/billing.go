package billing

import (
	"cmp"
	"context"
	"encoding/json"
	"log/slog"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/nexusmq/pkg/broker"
	"github.com/luponetn/noitrex/internal/db"
	"github.com/luponetn/noitrex/internal/plans"
)

type Billing struct {
	db     db.Querier
	broker broker.Broker
}

func NewBilling(db db.Querier, broker broker.Broker) *Billing {
	return &Billing{db: db, broker: broker}
}

func (b *Billing) Start(ctx context.Context) {
	slog.Info("Starting billing engine...")
	b.ProcessCustomerInvoices()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			slog.Info("Running scheduled billing process")
			b.ProcessCustomerInvoices()
		case <-ctx.Done():
			slog.Info("Billing engine shutting down")
			return
		}
	}
}

func (b *Billing) ProcessCustomerInvoices() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	customers, err := b.db.GetCustomersDueForBillingWithoutInvoice(ctx)
	if err != nil {
		slog.Error("failed to get customers due for billing without invoice", "error", err)
		return
	}

	groupedCustomers := b.GroupByCustomer(customers)

	for _, customer := range groupedCustomers {
		lineItems := make(map[string]any)
		events := []string{}
		totalUnitsUsed := int64(0)

		for _, aggregate := range customer.Aggregates {
			events = append(events, aggregate.EventName)
			totalUnitsUsed += aggregate.TotalUnits
		}

		var tiers []plans.PlanTiers
		if err := json.Unmarshal(customer.PlanTiers, &tiers); err != nil {
			slog.Error("failed to unmarshal plan tiers", "error", err)
			continue
		}

		slices.SortFunc(tiers, func(a, b plans.PlanTiers) int {
			return cmp.Compare(a.From, b.From)
		})

		remainingUnits := totalUnitsUsed
		var totalAmountInCents int64
		for _, tier := range tiers {

			if remainingUnits <= 0 {
				break
			}
			// Open-ended tier
			if tier.To == nil {
				totalAmountInCents += remainingUnits * tier.Price
				break
			}
			capacity := *tier.To - tier.From + 1
			billableUnits := min(remainingUnits, capacity)
			totalAmountInCents += billableUnits * tier.Price
			remainingUnits -= billableUnits
		}

		var newPeriodStart pgtype.Timestamptz
		newPeriodStart.Valid = true
		if customer.PlanBillingPeriod == db.BillingPeriodTypeYearly {
			newPeriodStart.Time = customer.PeriodStart.Time.AddDate(1, 0, 0)
		} else {
			newPeriodStart.Time = customer.PeriodStart.Time.AddDate(0, 1, 0)
		}

		lineItems["events"] = events
		lineItems["total_units"] = totalUnitsUsed
		lineItems["total_amount"] = totalAmountInCents
		lineItems["period_start"] = customer.PeriodStart.Time.Format(time.RFC3339)
		lineItems["period_end"] = newPeriodStart.Time.Format(time.RFC3339)
		lineItems["plan"] = customer.PlanName

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		lineItemsBytes, err := json.Marshal(lineItems)
		if err != nil {
			slog.Error("failed to marshal line items", "error", err)
			continue
		}

		newInvoice, err := b.db.CreateNewInvoice(ctx, db.CreateNewInvoiceParams{
			OperatorID:  customer.OperatorID,
			CustomerID:  customer.CustomerID,
			AmountCents: totalAmountInCents,
			Status:      "pending",
			PeriodStart: customer.PeriodStart,
			PeriodEnd:   newPeriodStart,
			LineItems:   lineItemsBytes,
		})

		err = b.db.UpdateCustomerPeriodStart(ctx, db.UpdateCustomerPeriodStartParams{
			ID:          customer.CustomerID,
			PeriodStart: newPeriodStart,
		})
		if err != nil {
			slog.Error("failed to update customer period start", "error", err)
		}

		msg, err := json.Marshal(newInvoice)
		if err != nil {
			slog.Error("failed to marshal invoice", "error", err)
			continue
		}

		b.broker.Publish("invoice.created", &broker.Message{
			Payload:   msg,
			Timestamp: time.Now(),
		})
	}

}

func (b *Billing) GroupByCustomer(customers []db.GetCustomersDueForBillingWithoutInvoiceRow) []CustomerForBilling {
	seen := make(map[string]int)
	var groupedCustomers []CustomerForBilling

	for _, c := range customers {
		existingIdx, ok := seen[c.CustomerID.String()]
		if ok {
			if c.EventName.Valid && c.TotalUnits.Valid {
				groupedCustomers[existingIdx].Aggregates = append(groupedCustomers[existingIdx].Aggregates, UsageAggregate{
					TotalUnits: c.TotalUnits.Int64,
					EventName:  c.EventName.String,
				})
			}
		} else {
			customer := CustomerForBilling{
				CustomerID:        c.CustomerID,
				OperatorID:        c.OperatorID,
				CustomerName:      c.CustomerName,
				PeriodStart:       c.PeriodStart,
				PlanName:          c.PlanName,
				PlanPricingModel:  c.PlanPricingModel,
				PlanTiers:         c.PlanTiers,
				PlanUnitPrice:     c.PlanUnitPrice,
				PlanBillingPeriod: c.PlanBillingPeriod,
			}

			if c.EventName.Valid && c.TotalUnits.Valid {
				customer.Aggregates = append(customer.Aggregates, UsageAggregate{
					TotalUnits: c.TotalUnits.Int64,
					EventName:  c.EventName.String,
				})
			}

			seen[c.CustomerID.String()] = len(groupedCustomers)
			groupedCustomers = append(groupedCustomers, customer)
		}
	}

	return groupedCustomers
}
