package customers

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

func (h *Handler) GetCustomer(c *gin.Context) {
	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		slog.Error("unauthorized access")
		utils.Unauthorized(c)
		return
	}

	var opID pgtype.UUID
	if err := opID.Scan(operatorIDStr.(string)); err != nil {
		slog.Error("Invalid operator_id format")
		utils.Unauthorized(c)
		return
	}

	customerIDStr := c.Param("id")
	var custID pgtype.UUID
	if err := custID.Scan(customerIDStr); err != nil {
		slog.Error("Invalid customer_id format", slog.String("id", customerIDStr))
		utils.BadRequest(c, "Invalid customer ID format")
		return
	}

	customer, err := h.service.GetCustomerByID(c.Request.Context(), custID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.NotFound(c, "Customer")
			return
		}
		slog.Error("failed to get customer", slog.Any("error", err))
		utils.InternalError(c)
		return
	}

	if customer.OperatorID != opID {
		utils.NotFound(c, "Customer")
		return
	}

	utils.OK(c, ToCustomerResponse(customer))
}

func (h *Handler) GetCustomerByEmail(c *gin.Context) {
	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		slog.Error("unauthorized request: operator id not found")
		utils.Unauthorized(c)
		return
	}

	var opID pgtype.UUID
	if err := opID.Scan(operatorIDStr.(string)); err != nil {
		slog.Error("something went wrong when trying to parse operator id")
		utils.InternalError(c)
		return
	}

	email := c.Query("email")
	if email == "" {
		slog.Error("email was not provided")
		utils.BadRequest(c, "Invalid request, email not provided.")
		return
	}

	customer, err := h.service.GetCustomerByEmail(c.Request.Context(), email, opID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.NotFound(c, "Customer")
			return
		}
		slog.Error("failed to get customer by email", slog.Any("error", err))
		utils.InternalError(c)
		return
	}

	utils.OK(c, ToCustomerResponse(customer))
}

func (h *Handler) GetCustomerByExternalID(c *gin.Context) {
	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		slog.Error("unauthorized request: operator id not found")
		utils.Unauthorized(c)
		return
	}

	var opID pgtype.UUID
	if err := opID.Scan(operatorIDStr.(string)); err != nil {
		slog.Error("something went wrong when trying to parse operator id")
		utils.InternalError(c)
		return
	}

	extIDStr := c.Query("external_id")
	if extIDStr == "" {
		slog.Error("external_id was not provided")
		utils.BadRequest(c, "Invalid request, external_id not provided.")
		return
	}

	var extID pgtype.UUID
	if err := extID.Scan(extIDStr); err != nil {
		slog.Error("Invalid external_id format", slog.String("external_id", extIDStr))
		utils.BadRequest(c, "Invalid external ID format")
		return
	}

	customer, err := h.service.GetCustomerByExternalID(c.Request.Context(), opID, extID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.NotFound(c, "Customer")
			return
		}
		slog.Error("failed to get customer by external ID", slog.Any("error", err))
		utils.InternalError(c)
		return
	}

	utils.OK(c, ToCustomerResponse(customer))
}

func (h *Handler) ListCustomers(c *gin.Context) {
	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		slog.Error("unauthorized access")
		utils.Unauthorized(c)
		return
	}

	var opID pgtype.UUID
	if err := opID.Scan(operatorIDStr.(string)); err != nil {
		slog.Error("Invalid operator_id format")
		utils.Unauthorized(c)
		return
	}

	customersList, err := h.service.ListCustomers(c.Request.Context(), opID)
	if err != nil {
		slog.Error("failed to list customers", slog.Any("error", err))
		utils.InternalError(c)
		return
	}

	responses := make([]*CustomerResponse, len(customersList))
	for i, customer := range customersList {
		responses[i] = ToCustomerResponse(customer)
	}

	utils.OK(c, responses)
}
