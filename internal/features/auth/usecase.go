package auth

import (
	"context"
	"go-service-boilerplate/internal/domain"
	"go-service-boilerplate/internal/features/user"
	"go-service-boilerplate/internal/utils"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Hasher defines methods for hashing and verifying passwords.
type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

// TokenGenerator defines methods for generating tokens.
type TokenGenerator interface {
	GenerateToken(userID uuid.UUID) (string, error)
}

type usecase struct {
	repo     Repository
	userRepo user.Repository
	hasher   Hasher
	jwt      TokenGenerator
	log      *zap.SugaredLogger
}

func NewUsecase(
	r Repository,
	ur user.Repository,
	h Hasher,
	j TokenGenerator,
	log *zap.SugaredLogger,
) Usecase {
	return &usecase{
		repo:     r,
		userRepo: ur,
		hasher:   h,
		jwt:      j,
		log:      log.Named("AuthUsecase"),
	}
}

func (u *usecase) Register(ctx context.Context, req RegisterRequest) (*LoginResponse, error) {
	// 1. Cek apakah email sudah ada
	existingUser, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		u.log.Error("failed to get user by email: ", err)
		return nil, domain.ErrInternal
	}

	if existingUser != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	// 2. Hash password
	hashedPass, err := u.hasher.HashPassword(req.Password)
	if err != nil {
		u.log.Error("failed to hash password: ", err)
		return nil, domain.ErrInternal
	}

	// Create username from email prefix
	username := strings.Split(req.Email, "@")[0]
	username = strings.ToLower(username) + utils.GenerateRandomString(6)

	// Save new user
	newUser := domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Username: username,
		Password: hashedPass,
	}

	createdUser, err := u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		u.log.Error("failed to create user", zap.Error(err))
		return nil, domain.ErrInternal
	}

	// Return login response
	return u.generateLoginResponse(createdUser)
}

func (u *usecase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		u.log.Error("failed to get user by email: ", err)
		return nil, domain.ErrInternal
	}

	// Check if user exists
	if user == nil {
		u.log.Warnf("user with email %s doesn't exists", req.Email)
		return nil, domain.ErrInvalidCredentials
	}

	// Verify password
	if !u.hasher.CheckPasswordHash(req.Password, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate login response
	return u.generateLoginResponse(user)
}

func (u *usecase) generateLoginResponse(user *domain.User) (*LoginResponse, error) {
	token, err := u.jwt.GenerateToken(user.ID)
	if err != nil {
		u.log.Error("failed to generate token", zap.Error(err))
		return nil, domain.ErrInternal
	}

	resUser := UserLoginResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}

	return &LoginResponse{User: resUser, AccessToken: token}, nil
}
