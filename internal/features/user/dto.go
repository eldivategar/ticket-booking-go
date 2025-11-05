package user

import "github.com/google/uuid"

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
}
