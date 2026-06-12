package webhook

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/internal/db"
)

type Service interface {
	ListDeliveries(ctx context.Context, operatorID pgtype.UUID, limit, offset int32) ([]db.WebhookDelivery, error)
	UpdateEndpointUrl(ctx context.Context, operatorID pgtype.UUID, newUrl string) (db.WebhookEndpoint, error)
	GetDeliveryStats(ctx context.Context, operatorID pgtype.UUID) (db.GetWebhookDeliveriesStatsRow, error)
}

type Svc struct {
	db db.Querier
}

func NewService(db db.Querier) Service {
	return &Svc{
		db: db,
	}
}

func (s *Svc) ListDeliveries(ctx context.Context, operatorID pgtype.UUID, limit, offset int32) ([]db.WebhookDelivery, error) {
	return s.db.ListWebhookDeliveriesPaginated(ctx, db.ListWebhookDeliveriesPaginatedParams{
		OperatorID: operatorID,
		Limit:      limit,
		Offset:     offset,
	})
}

func (s *Svc) UpdateEndpointUrl(ctx context.Context, operatorID pgtype.UUID, newUrl string) (db.WebhookEndpoint, error) {
	return s.db.UpdateWebhookEndpointURL(ctx, db.UpdateWebhookEndpointURLParams{
		OperatorID: operatorID,
		Url:        newUrl,
	})
}

func (s *Svc) GetDeliveryStats(ctx context.Context, operatorID pgtype.UUID) (db.GetWebhookDeliveriesStatsRow, error) {
	return s.db.GetWebhookDeliveriesStats(ctx, operatorID)
}