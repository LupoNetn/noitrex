package plans

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/internal/db"
	"github.com/luponetn/noitrex/utils"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) HandleCreatePlan(c *gin.Context) {
	var req CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Invalid request body", "error", err)
		utils.BadRequest(c, "Invalid request")
		return
	}

	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		slog.Error("unauthorized access, operator id not found")
		utils.Unauthorized(c)
		return
	}

	var opId pgtype.UUID
	err := opId.Scan(operatorIDStr)
	if err != nil {
		slog.Error("failed to parse operator id", "error", err)
		utils.BadRequest(c, "Invalid operator id")
		return
	}

	params := db.CreatePlanParams{
		Name:          req.Name,
		PricingModel:  req.PricingModel,
		UnitPriceCent: req.UnitPriceCent,
		Tiers:         req.Tiers,
		BillingPeriod: req.BillingPeriod,
		OperatorID:    opId,
	}

	plan, err := h.service.CreatePlan(c.Request.Context(), params)
	if err != nil {
		if errors.Is(err, ErrPlanAlreadyExists) {
			slog.Info("plan already exists")
			utils.OK(c, plan)
		}
		slog.Error("failed to create plan", "error", err)
		utils.InternalError(c)
		return
	}

	utils.OK(c, plan)
}

func (h *Handler) HandleGetPlanById(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		slog.Error("missing plan id")
		utils.BadRequest(c, "missing plan id")
		return
	}

	var id pgtype.UUID
	err := id.Scan(idStr)
	if err != nil {
		slog.Error("failed to parse plan id", "error", err)
		utils.BadRequest(c, "invalid plan id")
		return
	}

	plan, err := h.service.GetPlanById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPlanNotFound) {
			slog.Warn("plan not found")
			utils.NotFound(c, "plan not found")
			return
		}
		slog.Error("failed to get plan", "error", err)
		utils.InternalError(c)
		return
	}

	utils.OK(c, plan)
}

func (h *Handler) HandleGetPlanByName(c *gin.Context) {
	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		slog.Error("unauthorized access, operator id not found")
		utils.Unauthorized(c)
		return
	}

	var opId pgtype.UUID
	err := opId.Scan(operatorIDStr)
	if err != nil {
		slog.Error("failed to parse operator id", "error", err)
		utils.BadRequest(c, "Invalid operator id")
		return
	}

	plan, err := h.service.GetPlanByName(c.Request.Context(), db.GetPlanByNameParams{
		OperatorID: opId,
		Name:       c.Query("name"),
	})
	if err != nil {
		if errors.Is(err, ErrPlanNotFound) {
			slog.Warn("plan not found")
			utils.NotFound(c, "plan not found")
			return
		}
		slog.Error("failed to get plan", "error", err)
		utils.InternalError(c)
		return
	}

	utils.OK(c, plan)
}

func (h *Handler) HandleListPlans(c *gin.Context) {
	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		slog.Error("unauthorized access, operator id not found")
		utils.Unauthorized(c)
		return
	}

	var opId pgtype.UUID
	err := opId.Scan(operatorIDStr)
	if err != nil {
		slog.Error("failed to parse operator id", "error", err)
		utils.BadRequest(c, "Invalid operator id")
		return
	}

	plans, err := h.service.ListPlans(c.Request.Context(), opId)
	if err != nil {
		slog.Error("failed to list plans", "error", err)
		utils.InternalError(c)
		return
	}

	utils.OK(c, plans)
}

func (h *Handler) HandleUpdatePlan(c *gin.Context) {
	var req UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Invalid request body", "error", err)
		utils.BadRequest(c, "Invalid request")
		return
	}

	params := db.UpdatePlanParams{
		ID:            req.ID,
		Name:          req.Name,
		PricingModel:  req.PricingModel,
		UnitPriceCent: req.UnitPriceCent,
		Tiers:         req.Tiers,
		BillingPeriod: req.BillingPeriod,
	}

	plan, err := h.service.UpdatePlan(c.Request.Context(), params)
	if err != nil {
		if errors.Is(err, ErrPlanAlreadyExists) {
			slog.Info("plan already exists")
			utils.OK(c, plan)
		}
		slog.Error("failed to update plan", "error", err)
		utils.InternalError(c)
		return
	}

	utils.OK(c, plan)
}

func (h *Handler) HandleDeletePlan(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		slog.Error("missing plan id")
		utils.BadRequest(c, "missing plan id")
		return
	}

	var id pgtype.UUID
	err := id.Scan(idStr)
	if err != nil {
		slog.Error("failed to parse plan id", "error", err)
		utils.BadRequest(c, "invalid plan id")
		return
	}

	err = h.service.DeletePlan(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPlanNotFound) {
			slog.Warn("plan not found")
			utils.NotFound(c, "plan not found")
			return
		}
		slog.Error("failed to delete plan", "error", err)
		utils.InternalError(c)
		return
	}

	utils.OK(c, nil)
}
