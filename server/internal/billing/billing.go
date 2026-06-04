package billing

import (
	"context"
	"log/slog"

	"github.com/luponetn/nexusmq/pkg/broker"
	"github.com/luponetn/noitrex/internal/db"
)

type Billing struct {
	db     db.Querier
	broker broker.Broker
}

func NewBilling(db db.Querier, broker broker.Broker) *Billing {
	return &Billing{db: db, broker: broker}
}

func (b *Billing) Start() {}

func (b *Billing) ProcessCustomerInvoices() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	customers, err := b.db.GetCustomersDueForBillingWithoutInvoice(ctx)
	if err != nil {
		slog.Error("failed to get customers due for billing without invoice", "error", err)
		return
	}

	groupedCustomers := b.GroupByCustomer(customers)

}

func (b *Billing) GroupByCustomer(customers []db.GetCustomersDueForBillingWithoutInvoiceRow) []CustomerForBilling {
	seen := make(map[string]int)
	var groupedCustomers []CustomerForBilling

	for idx, c := range customers {
		_, ok := seen[c.CustomerID.String()]
		if ok {
			if c.EventName.Valid && c.TotalUnits.Valid {
				groupedCustomers[idx].Aggregates = append(groupedCustomers[idx].Aggregates, UsageAggregate{
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

			groupedCustomers = append(groupedCustomers, customer)
			seen[c.CustomerID.String()] = idx
		}
	}

	return groupedCustomers
}
