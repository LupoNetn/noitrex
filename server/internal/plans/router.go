package plans

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/luponetn/noitrex/internal/middlewares"
)

func NewRouter(router *gin.Engine, h *Handler, JWTAccessSecret string) {
	planGroup := router.Group("/plans")

	planGroup.Use(middleware.AuthMiddleware(JWTAccessSecret))

	planGroup.POST("/", h.HandleCreatePlan)
	planGroup.GET("/:id", h.HandleGetPlanById)
	planGroup.GET("/", h.HandleGetPlanByName)
	planGroup.GET("/list", h.HandleListPlans)
	planGroup.PATCH("/:id", h.HandleUpdatePlan)
	planGroup.DELETE("/:id", h.HandleDeletePlan)
}
