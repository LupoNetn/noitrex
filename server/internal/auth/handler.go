package auth

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func NewHandler(Svc Service) *Handler {
	return &Handler{
		service: Svc,
	}
}

// implement all http handlers
func (h *Handler) HandleLogin(c *gin.Context) {
	
}

func (h *Handler) HandleRegister(c *gin.Context) {
	
}
