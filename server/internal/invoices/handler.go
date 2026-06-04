package invoices

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) HandleGenerateCustomerInvoice(c *gin.Context) {}
