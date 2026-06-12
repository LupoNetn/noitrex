package webhook

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/luponetn/noitrex/internal/middlewares"
)

func NewRouter(router *gin.Engine, h *Handler, jwtSecret string) {
	webhookGroup := router.Group("/webhooks")
	webhookGroup.Use(middleware.AuthMiddleware(jwtSecret))

	webhookGroup.GET("/deliveries", h.HandleGetWebhookHistory)
	webhookGroup.PUT("/endpoint/url", h.HandleUpdateEndpointUrl)
	webhookGroup.GET("/deliveries/stats", h.HandleGetDeliveriesStats)
}