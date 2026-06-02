package customers

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

func NewHandler(svc Service) *Handler {
	return &Handler{
		service: svc,
	}
}

func (h *Handler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind customer request", slog.Any("error", err))
		utils.BadRequest(c, "Invalid request body")
		return
	}

	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		slog.Error("unathorized access")
		utils.Unauthorized(c)
		return
	}

	var opID pgtype.UUID
	if err := opID.Scan(operatorIDStr.(string)); err != nil {
		slog.Error("Invalid operator_id format")
		utils.Unauthorized(c)
		return
	}
	var exID pgtype.UUID
	if err := exID.Scan(req.ExternalID); err != nil {
		slog.Error("Invalid external_id format")
		utils.BadRequest(c, "Invalid external_id format")
		return
	}

	params := db.CreateCustomerParams{
		OperatorID: opID,
		ExternalID: exID,
		Name:       req.Name,
		Email:      req.Email,
	}

	customer, err := h.service.CreateCustomer(c.Request.Context(), params)
	if err != nil {
		var customerAlreadyExistsErr *CustomerAlreadyExists
		if errors.As(err, &customerAlreadyExistsErr) {
			slog.Warn("customer already exists", slog.String("customer_id", customerAlreadyExistsErr.ExternalID), slog.String("operator_id", customerAlreadyExistsErr.OperatorID))
			utils.Conflict(c, customerAlreadyExistsErr.Error())
			return
		}
		slog.Error("failed to create customer", slog.Any("error", err))
		utils.InternalError(c)
		return
	}
	utils.Created(c, ToCustomerResponse(customer))
}
