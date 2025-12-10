package app

import (
	"go-war-ticket-service/configs"
	"go-war-ticket-service/internal/features/auth"
	"go-war-ticket-service/internal/features/event"
	"go-war-ticket-service/internal/features/order"
	"go-war-ticket-service/internal/features/ticket"
	"go-war-ticket-service/internal/features/user"
	"go-war-ticket-service/internal/platform/hash"
	"go-war-ticket-service/internal/platform/jwt"
	rabbitmq "go-war-ticket-service/internal/platform/message_broker/rabbit_mq"
	"go-war-ticket-service/internal/platform/middleware"
	"go-war-ticket-service/internal/platform/pdf"
	"go-war-ticket-service/internal/platform/validator"
	"go-war-ticket-service/internal/utils"

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
	EventHandler   event.Handler
	OrderHandler   order.Handler
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
	authMiddleware := middleware.AuthRequired(cfg.JWTAccessSecret, log)
	mqPublisher, err := rabbitmq.NewRabbitMQPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Error("Failed to create RabbitMQ publisher", zap.Error(err))
		return nil
	}
	pdfGenerator := pdf.NewMarotoGenerator()

	// Create new queue
	mqPublisher.CreateQueue(utils.QueueTicketGeneration)

	// User Features
	userRepo := user.NewRepository(db)
	userUsecase := user.NewUsecase(userRepo, log, s3, cfg)
	userHandler := user.NewHandler(userUsecase)

	// Auth Features
	authRepo := auth.NewRepository(db)
	authUsecase := auth.NewUsecase(authRepo, userRepo, hasher, jwtGen, log, s3, cfg)
	authHandler := auth.NewHandler(authUsecase, val)

	// Event Features
	eventRepo := event.NewRepository(db)
	eventUsecase := event.NewUsecase(eventRepo, log, s3, cfg)
	eventHandler := event.NewHandler(eventUsecase, val)

	// Order Features
	orderRepo := order.NewRepository(db)
	orderUsecase := order.NewUsecase(orderRepo, log, s3, cfg, rdb)
	orderService := order.NewService(orderRepo, log, mqPublisher)
	orderHandler := order.NewHandler(orderUsecase, orderService, val)

	// Worker
	ticketRepo := ticket.NewRepository(db)
	ticketWorker := ticket.NewTicketWorker(mqPublisher.GetConnection(), ticketRepo, orderRepo, s3, cfg, pdfGenerator, log)

	go ticketWorker.Start()

	return &Dependencies{
		AuthHandler:    *authHandler,
		UserHandler:    *userHandler,
		AuthMiddleware: authMiddleware,
		EventHandler:   *eventHandler,
		OrderHandler:   *orderHandler,
	}
}
