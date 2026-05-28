package operator

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
func (h *Handler) HandleGetOperatorByID(c *gin.Context) {
	
}