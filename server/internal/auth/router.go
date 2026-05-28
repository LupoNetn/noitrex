package auth

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, h *Handler) {
	authGroup := router.Group("/auth")

	authGroup.POST("/login", h.HandleLogin)
	authGroup.POST("/register", h.HandleRegister)
}
