package events

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/internal/db"
	"github.com/luponetn/noitrex/utils"
)

type Handler struct {
	Service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) HandleCreateEvent(c *gin.Context) {
	var req UsageEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Invalid request payload")
		utils.BadRequest(c, "Invalid request payload")
		return
	}

	operatorID, exists := c.Get("operator_id")
	if !exists {
		slog.Error("unathorized access")
		utils.Unauthorized(c)
		return
	}

	var customerID pgtype.UUID
	err := customerID.Scan(req.CustomerID)
	if err != nil {
		slog.Error("Invalid customer_id format")
		utils.BadRequest(c, "Invalid customer_id format")
		return
	}

	params := db.CreateUsageEventParams{
		CustomerID:     customerID,
		OperatorID:     operatorID.(pgtype.UUID),
		EventName:      req.EventName,
		Quantity:       req.Quantity,
		IdempotencyKey: req.IdempotencyKey,
	}

	event, err := h.Service.CreateUsageEvent(c.Request.Context(), params)
	if err != nil {
		var u *UsageEventExists
		if errors.As(err, &u) {
			slog.Error("Usage event already exists")
			utils.BadRequest(c, u.Error())
			return
		}
		slog.Error("Error creating usage event")
		utils.InternalError(c)
		return
	}

	utils.Created(c, event)
}
