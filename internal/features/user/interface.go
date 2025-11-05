package user

import (
	"context"
	"go-service-boilerplate/internal/domain"
)

type Usecase interface{}

type Repository interface {
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
}
