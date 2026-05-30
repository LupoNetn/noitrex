package auth

import "errors"

// Errors(Sentinel and Custom Types)
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// Request and Response Structures
type RegisterRequest struct {
	Name         string  `json:"name" binding:"required"`
	LogoUrl      *string `json:"logo_url" binding:"omitempty,url"`
	Description  *string `json:"description" binding:"omitempty"`
	SupportEmail string  `json:"support_email" binding:"required,email"`
	WebsiteUrl   *string `json:"website_url" binding:"omitempty,url"`
	Password     string  `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type OperatorDTO struct {
	ID           any     `json:"id"`
	Name         string  `json:"name"`
	LogoUrl      *string `json:"logo_url,omitempty"`
	Description  *string `json:"description,omitempty"`
	SupportEmail string  `json:"support_email"`
	WebsiteUrl   *string `json:"website_url,omitempty"`
	Status       string  `json:"status"`
}
