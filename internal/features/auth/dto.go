package auth

import "github.com/google/uuid"

type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	// Username string `json:"username" validate:"required,alphanum,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Avatar   string `json:"avatar" validate:"omitempty,base64"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Avatar   string    `json:"avatar"`
}

type LoginResponse struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	User         UserLoginResponse `json:"user"`
}
