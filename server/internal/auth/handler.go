package auth

import (
	"log/slog"
	"errors"

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
		if errors.Is(err, ErrInvalidCredentials) {
			utils.Unauthorized(c)
			return
		}
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
		if errors.Is(err, ErrUserAlreadyExists) {
			utils.Conflict(c, "User with this email already exists")
			return
		}
		slog.Error("Error creating operator", "error", err.Error())
		utils.InternalError(c)
		return
	}

	dto := OperatorDTO{
		ID:           operator.ID,
		Name:         operator.Name,
		LogoUrl:      nil,
		Description:  nil,
		SupportEmail: operator.SupportEmail.String,
		WebsiteUrl:   nil,
		Status:       string(operator.Status),
	}
	if operator.LogoUrl.Valid {
		dto.LogoUrl = &operator.LogoUrl.String
	}
	if operator.Description.Valid {
		dto.Description = &operator.Description.String
	}
	if operator.WebsiteUrl.Valid {
		dto.WebsiteUrl = &operator.WebsiteUrl.String
	}

	utils.Created(c, gin.H{
		"operator":       dto,
		"api_key":        apiKey,
		"webhook_secret": webhookSecret,
	})
}
