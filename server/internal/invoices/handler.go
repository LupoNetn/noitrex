package invoices

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/utils"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetInvoice(c *gin.Context) {
	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		utils.Unauthorized(c)
		return
	}
	var opID pgtype.UUID
	if err := opID.Scan(operatorIDStr.(string)); err != nil {
		utils.Unauthorized(c)
		return
	}

	invoiceIDStr := c.Param("id")
	var invID pgtype.UUID
	if err := invID.Scan(invoiceIDStr); err != nil {
		utils.BadRequest(c, "Invalid invoice ID format")
		return
	}

	inv, err := h.service.GetInvoiceByID(c.Request.Context(), opID, invID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, ErrUnauthorized) {
			utils.NotFound(c, "Invoice")
			return
		}
		slog.Error("failed to get invoice", slog.Any("error", err))
		utils.InternalError(c)
		return
	}

	utils.OK(c, ToInvoiceResponse(inv))
}

func (h *Handler) ListCustomerInvoices(c *gin.Context) {
	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		utils.Unauthorized(c)
		return
	}
	var opID pgtype.UUID
	if err := opID.Scan(operatorIDStr.(string)); err != nil {
		utils.Unauthorized(c)
		return
	}

	customerIDStr := c.Param("customerId")
	if customerIDStr == "" {
		customerIDStr = c.Query("customer_id")
	}
	var custID pgtype.UUID
	if err := custID.Scan(customerIDStr); err != nil {
		utils.BadRequest(c, "Invalid customer ID format")
		return
	}

	invoicesList, err := h.service.ListInvoicesByCustomer(c.Request.Context(), opID, custID)
	if err != nil {
		slog.Error("failed to list invoices", slog.Any("error", err))
		utils.InternalError(c)
		return
	}

	responses := make([]*InvoiceResponse, len(invoicesList))
	for i, inv := range invoicesList {
		responses[i] = ToInvoiceResponse(inv)
	}

	utils.OK(c, responses)
}

func (h *Handler) UpdateInvoice(c *gin.Context) {
	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		utils.Unauthorized(c)
		return
	}
	var opID pgtype.UUID
	if err := opID.Scan(operatorIDStr.(string)); err != nil {
		utils.Unauthorized(c)
		return
	}

	invoiceIDStr := c.Param("id")
	var invID pgtype.UUID
	if err := invID.Scan(invoiceIDStr); err != nil {
		utils.BadRequest(c, "Invalid invoice ID format")
		return
	}

	var req UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	inv, err := h.service.UpdateInvoiceStatus(c.Request.Context(), opID, invID, req.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, ErrUnauthorized) {
			utils.NotFound(c, "Invoice")
			return
		}
		slog.Error("failed to update invoice", slog.Any("error", err))
		utils.InternalError(c)
		return
	}

	utils.OK(c, ToInvoiceResponse(inv))
}

func (h *Handler) ListOperatorInvoices(c *gin.Context) {
	operatorIDStr, exists := c.Get("operatorId")
	if !exists {
		utils.Unauthorized(c)
		return
	}
	var opID pgtype.UUID
	if err := opID.Scan(operatorIDStr.(string)); err != nil {
		utils.Unauthorized(c)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	invList, total, err := h.service.ListOperatorInvoices(c.Request.Context(), opID, int32(limit), int32(offset))
	if err != nil {
		slog.Error("failed to list operator invoices", slog.Any("error", err))
		utils.InternalError(c)
		return
	}

	responses := make([]*InvoiceResponse, len(invList))
	for i, inv := range invList {
		responses[i] = ToInvoiceResponse(inv)
	}

	utils.OK(c, PaginatedInvoicesResponse{
		Invoices: responses,
		Page:     page,
		Limit:    limit,
		Total:    total,
	})
}
