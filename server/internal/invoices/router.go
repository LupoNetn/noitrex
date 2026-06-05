package invoices

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/luponetn/noitrex/internal/middlewares"
)

func NewRouter(router *gin.Engine, h *Handler, JWTAccessSecret string) {
	invoicesGroup := router.Group("/invoices")
	invoicesGroup.Use(middleware.AuthMiddleware(JWTAccessSecret))

	invoicesGroup.GET("", h.ListOperatorInvoices)         // GET /invoices?page=1&limit=20
	invoicesGroup.GET("/:id", h.GetInvoice)               // GET /invoices/:id
	invoicesGroup.GET("/customer", h.ListCustomerInvoices) // GET /invoices/customer?customer_id=<uuid>
	invoicesGroup.PATCH("/:id", h.UpdateInvoice)          // PATCH /invoices/:id
}
