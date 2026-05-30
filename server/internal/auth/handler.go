package auth

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/luponetn/noitrex/utils"
)

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
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Bad request when logging in operator", "error", err.Error())
		utils.BadRequest(c, "Invalid request payload")
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		slog.Error("Error logging in operator", "error", err.Error())
		utils.InternalError(c)
		return
	}

	utils.OK(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) HandleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Bad request when creating operator", "error", err.Error())
		utils.BadRequest(c, "Invalid request payload")
		return
	}

	operator, apiKey, webhookSecret, err := h.service.CreateOperator(c.Request.Context(), req)
	if err != nil {
		slog.Error("Error creating operator", "error", err.Error())
		utils.InternalError(c)
		return
	}

	operator.PasswordHash.String = ""
	operator.ApiKeyHash = ""
	operator.WebhookSecret = ""

	utils.Created(c, gin.H{
		"operator":       operator,
		"api_key":        apiKey,
		"webhook_secret": webhookSecret,
	})
}
