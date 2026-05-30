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

	// Set HttpOnly Cookies
	// maxAge for access token: 15 minutes = 900 seconds
	// maxAge for refresh token: 7 days = 604800 seconds
	// Note: Set 'secure' (second to last boolean) to true in production over HTTPS
	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/auth/refresh", "", false, true)

	utils.OK(c, gin.H{
		"message": "logged in successfully",
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

func (h *Handler) HandleRefresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		slog.Error("Refresh token cookie not found", "error", err.Error())
		utils.Unauthorized(c)
		return
	}

	accessToken, newRefreshToken, err := h.service.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			utils.Unauthorized(c)
			return
		}
		slog.Error("Error refreshing token", "error", err.Error())
		utils.InternalError(c)
		return
	}

	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", newRefreshToken, 604800, "/auth/refresh", "", false, true)

	utils.OK(c, gin.H{
		"message": "tokens refreshed successfully",
	})
}

func (h *Handler) HandleLogout(c *gin.Context) {
	// Clear the HttpOnly cookies by setting MaxAge to -1
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/auth/refresh", "", false, true)

	utils.OK(c, gin.H{
		"message": "logged out successfully",
	})
}
