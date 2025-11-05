package auth

import (
	"context"
)

type Usecase interface {
	Register(ctx context.Context, req RegisterRequest) (*LoginResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
}

type Repository interface{}
