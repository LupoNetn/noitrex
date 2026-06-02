package customers

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/luponetn/noitrex/internal/middlewares"
)

func NewRouter(router *gin.Engine, h *Handler, JWTAccessSecret string) {
	customerGroup := router.Group("/customers")
	customerGroup.Use(middleware.AuthMiddleware(JWTAccessSecret))

	customerGroup.POST("/", h.CreateCustomer)
	customerGroup.GET("/", h.ListCustomers)
	customerGroup.GET("/email", h.GetCustomerByEmail)
	customerGroup.GET("/external", h.GetCustomerByExternalID)
	customerGroup.GET("/:id", h.GetCustomer)
}
