package user

import (
	"context"
	"go-war-ticket-service/internal/domain"

	"github.com/google/uuid"
)

type Usecase interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

type Repository interface {
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}
