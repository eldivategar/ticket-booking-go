package app

import (
	"fmt"
	"go-service-boilerplate/configs"
	"go-service-boilerplate/internal/platform/cache"
	"go-service-boilerplate/internal/platform/database"
	"go-service-boilerplate/internal/platform/storage"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	app *fiber.App
	cfg configs.Config
	log *zap.SugaredLogger
}

func NewServer(cfg configs.Config, log *zap.SugaredLogger) *Server {
	app := fiber.New()

	app.Use(recover.New())
	// app.Use(logger.New()) // fiber logger middleware

	return &Server{
		app: app,
		cfg: cfg,
		log: log,
	}
}

func (s *Server) Start() error {
	// Database connection
	db, err := database.Connect(s.cfg)
	if err != nil {
		s.log.Fatal("failed to connect to database", zap.Error(err))
	}
	s.log.Info("database connection established")

	// Redis connection
	rdb, err := cache.NewRedisClient(s.cfg)
	if err != nil {
		s.log.Fatal("failed to connect to redis", zap.Error(err))
	}
	s.log.Info("redis connection established")

	// Minio connection
	s3Client, err := storage.NewMinIOClient(s.cfg)
	if err != nil {
		s.log.Fatal("failed to connect to minio", zap.Error(err))
	}
	s.log.Info("minio (s3) connection established")

	// Setup Dependencies
	deps := SetupDependencies(s.cfg, s.log, db, rdb, s3Client)

	// Setup Routes
	SetupRoutes(s.app, deps)

	// Start server
	addr := fmt.Sprintf("%s:%d", s.cfg.ServerHost, s.cfg.ServerPort)
	s.log.Infof("starting server on %s", addr)
	return s.app.Listen(addr)
}
