package auth

import (
	"context"
	"go-service-boilerplate/configs"
	"go-service-boilerplate/internal/domain"
	"go-service-boilerplate/internal/features/user"
	"go-service-boilerplate/internal/platform/storage"
	"go-service-boilerplate/internal/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
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
	repo        Repository
	userRepo    user.Repository
	hasher      Hasher
	jwt         TokenGenerator
	log         *zap.SugaredLogger
	minioClient *minio.Client
	cfg         configs.Config
}

func NewUsecase(
	r Repository,
	ur user.Repository,
	h Hasher,
	j TokenGenerator,
	log *zap.SugaredLogger,
	minioClient *minio.Client,
	cfg configs.Config,
) Usecase {
	return &usecase{
		repo:        r,
		userRepo:    ur,
		hasher:      h,
		jwt:         j,
		log:         log.Named("AuthUsecase"),
		minioClient: minioClient,
		cfg:         cfg,
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
	username = strings.ToLower(username) + "-" + utils.GenerateRandomNumberString(6)

	// Save avatar to MinIO if provided
	if req.Avatar != "" {
		path, err := storage.UploadImageToMinIO(
			u.minioClient,
			u.cfg.MinioBucket,
			req.Avatar,
			"avatars",
			username,
		)
		if err != nil {
			u.log.Error("failed to upload avatar to MinIO: ", err)
			return nil, domain.ErrInternal
		}

		req.Avatar = path
	}

	// Save new user
	newUser := domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Username: username,
		Password: hashedPass,
		Avatar:   req.Avatar,
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

	// Create presigned url for avatar if exists
	avatarUrl := ""
	if user.Avatar != "" {
		presignedURL, err := storage.GetPresignedObject(
			u.minioClient,
			u.cfg.MinioBucket,
			user.Avatar,
			u.cfg.MinioEndpoint,
			u.cfg.MinioEndpoint,
			time.Minute*15,
		)
		if err != nil {
			u.log.Error("failed to generate presigned url for avatar: ", err)
			return nil, domain.ErrInternal
		}

		avatarUrl = presignedURL
	}

	resUser := UserLoginResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Avatar:   avatarUrl,
	}

	return &LoginResponse{User: resUser, AccessToken: token}, nil
}
