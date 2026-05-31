package events

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/luponetn/noitrex/internal/middlewares"
)

func NewRouter(router *gin.Engine, h *Handler, jwtAccessSecret string) {
	eventsGroup := router.Group("/events")
	eventsGroup.Use(middleware.AuthMiddleware(jwtAccessSecret))

	eventsGroup.POST("/", h.HandleCreateEvent)
}
