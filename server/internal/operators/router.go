package operator

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, h *Handler) {
	operatorGroup := router.Group("/operator")

	operatorGroup.GET("/:id", h.HandleGetOperatorByID)
}
