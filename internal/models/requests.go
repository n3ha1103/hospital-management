package models

// LoginRequest represents the login request payload

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin doctor nurse staff patient"`
	Phone    string `json:"phone,omitempty"`
}
