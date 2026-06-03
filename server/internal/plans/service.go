package plans

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/internal/db"
)

type Service interface {
	CreatePlan(ctx context.Context, args db.CreatePlanParams) (db.Plan, error)
	GetPlanById(ctx context.Context, id pgtype.UUID) (db.Plan, error)
	GetPlanByName(ctx context.Context, args db.GetPlanByNameParams) (db.Plan, error)
	ListPlans(ctx context.Context, operatorID pgtype.UUID) ([]db.Plan, error)
	UpdatePlan(ctx context.Context, args db.UpdatePlanParams) (db.Plan, error)
	DeletePlan(ctx context.Context, id pgtype.UUID) error
}

type Svc struct {
	db db.Querier
}

func NewService(db db.Querier) Service {
	return &Svc{db: db}
}

func (s *Svc) CreatePlan(ctx context.Context, args db.CreatePlanParams) (db.Plan, error) {
	existingPlan, err := s.db.GetPlanByName(ctx, db.GetPlanByNameParams{
		OperatorID: args.OperatorID,
		Name:       args.Name,
	})
	if err == nil {
		slog.Warn("Plan already exists", slog.String("plan_name", args.Name))
		return existingPlan, ErrPlanAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Failed to get plan by name", slog.String("plan_name", args.Name), slog.String("operator_id", args.OperatorID.String()))
		return db.Plan{}, err
	}

	return s.db.CreatePlan(ctx, args)
}

func (s *Svc) GetPlanById(ctx context.Context, id pgtype.UUID) (db.Plan, error) {
	plan, err := s.db.GetPlan(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("Plan not found", slog.String("plan_id", id.String()))
			return db.Plan{}, ErrPlanNotFound
		}
		slog.Error("Failed to get plan", slog.String("plan_id", id.String()), slog.String("error", err.Error()))
		return db.Plan{}, err
	}
	return plan, nil
}

func (s *Svc) GetPlanByName(ctx context.Context, args db.GetPlanByNameParams) (db.Plan, error) {
	plan, err := s.db.GetPlanByName(ctx, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("Plan not found", slog.String("plan_name", args.Name))
			return db.Plan{}, ErrPlanNotFound
		}
		slog.Error("Failed to get plan", slog.String("plan_name", args.Name), slog.String("operator_id", args.OperatorID.String()))
		return db.Plan{}, err
	}
	return plan, nil
}

func (s *Svc) ListPlans(ctx context.Context, operatorID pgtype.UUID) ([]db.Plan, error) {
	plans, err := s.db.ListPlans(ctx, operatorID)
	if err != nil {
		slog.Error("Failed to list plans", slog.String("operator_id", operatorID.String()), slog.String("error", err.Error()))
		return nil, err
	}
	return plans, nil
}

func (s *Svc) UpdatePlan(ctx context.Context, args db.UpdatePlanParams) (db.Plan, error) {
	plan, err := s.db.UpdatePlan(ctx, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("Plan not found", slog.String("plan_id", args.ID.String()))
			return db.Plan{}, ErrPlanNotFound
		}
		slog.Error("Failed to update plan", slog.String("plan_id", args.ID.String()), slog.String("error", err.Error()))
		return db.Plan{}, err
	}
	return plan, nil
}

func (s *Svc) DeletePlan(ctx context.Context, id pgtype.UUID) error {
	err := s.db.DeletePlan(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("Plan not found", slog.String("plan_id", id.String()))
			return ErrPlanNotFound
		}
		slog.Error("Failed to delete plan", slog.String("plan_id", id.String()), slog.String("error", err.Error()))
		return err
	}
	return nil
}
