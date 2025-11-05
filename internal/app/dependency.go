package app

import (
	"go-service-boilerplate/configs"
	"go-service-boilerplate/internal/features/auth"
	"go-service-boilerplate/internal/features/user"
	"go-service-boilerplate/internal/platform/hash"
	"go-service-boilerplate/internal/platform/jwt"
	"go-service-boilerplate/internal/platform/middleware"
	"go-service-boilerplate/internal/platform/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Dependencies holds all the dependencies for the application
type Dependencies struct {
	AuthHandler    auth.Handler
	UserHandler    user.Handler
	AuthMiddleware fiber.Handler
}

// Initialize and set up all dependencies
func SetupDependencies(
	cfg configs.Config,
	log *zap.SugaredLogger,
	db *gorm.DB,
	rdb *redis.Client,
	s3 *minio.Client,
) *Dependencies {
	// Platform
	hasher := hash.NewBcryptHasher()
	jwtGen := jwt.NewJWTGenerator(cfg.JWTAccessSecret)
	val := validator.New()
	authMiddleware := middleware.AuthRequired(cfg.JWTAccessSecret)

	// User Features
	userRepo := user.NewRepository(db)
	userUsecase := user.NewUsecase(userRepo, log)
	userHandler := user.NewHandler(userUsecase)

	// Auth Features
	authRepo := auth.NewRepository(db)
	authUsecase := auth.NewUsecase(authRepo, userRepo, hasher, jwtGen, log)
	authHandler := auth.NewHandler(authUsecase, val)

	return &Dependencies{
		AuthHandler:    *authHandler,
		UserHandler:    *userHandler,
		AuthMiddleware: authMiddleware,
	}
}
