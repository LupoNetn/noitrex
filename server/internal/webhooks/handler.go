package webhook

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/utils"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) HandleGetWebhookHistory(c *gin.Context) {
	var req PaginatedRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Error("Invalid pagination request", "error", err)
		utils.BadRequest(c, "Invalid pagination parameters")
		return
	}

	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		slog.Error("unauthorized access, operator id not found")
		utils.Unauthorized(c)
		return
	}

	var opId pgtype.UUID
	if err := opId.Scan(operatorIDStr); err != nil {
		slog.Error("failed to parse operator id", "error", err)
		utils.BadRequest(c, "Invalid operator id")
		return
	}

	deliveries, err := h.service.ListDeliveries(c.Request.Context(), opId, req.Limit, req.Offset)
	if err != nil {
		slog.Error("failed to list webhook deliveries", "error", err)
		utils.InternalError(c)
		return
	}

	utils.OK(c, PaginatedResponse{
		Data:   deliveries,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
}

func (h *Handler) HandleUpdateEndpointUrl(c *gin.Context) {
	var req UpdateEndpointUrlRequest
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
	if err := opId.Scan(operatorIDStr); err != nil {
		slog.Error("failed to parse operator id", "error", err)
		utils.BadRequest(c, "Invalid operator id")
		return
	}

	endpoint, err := h.service.UpdateEndpointUrl(c.Request.Context(), opId, req.URL)
	if err != nil {
		slog.Error("failed to update webhook endpoint url", "error", err)
		utils.InternalError(c)
		return
	}

	utils.OK(c, endpoint)
}

func (h *Handler) HandleGetDeliveriesStats(c *gin.Context) {
	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		slog.Error("unauthorized access, operator id not found")
		utils.Unauthorized(c)
		return
	}

	var opId pgtype.UUID
	if err := opId.Scan(operatorIDStr); err != nil {
		slog.Error("failed to parse operator id", "error", err)
		utils.BadRequest(c, "Invalid operator id")
		return
	}

	stats, err := h.service.GetDeliveryStats(c.Request.Context(), opId)
	if err != nil {
		slog.Error("failed to get webhook deliveries stats", "error", err)
		utils.InternalError(c)
		return
	}

	utils.OK(c, WebhookDeliveryStatsResponse{
		TotalDeliveries:      stats.TotalDeliveries,
		SuccessfulDeliveries: stats.SuccessfulDeliveries,
		FailedDeliveries:     stats.FailedDeliveries,
	})
}
