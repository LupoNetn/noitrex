package auth

import "errors"

// Errors(Sentinel and Custom Types)
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Request and Response Structures
type RegisterRequest struct {
	Name         string  `json:"name" validate:"required"`
	LogoUrl      *string `json:"logo_url" validate:"omitempty,url"`
	Description  *string `json:"description" validate:"omitempty"`
	SupportEmail string  `json:"support_email" validate:"required,email"`
	WebsiteUrl   *string `json:"website_url" validate:"omitempty,url"`
	Password     string  `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
