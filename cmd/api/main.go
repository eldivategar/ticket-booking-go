package main

import (
	"go-war-ticket-service/configs"
	"go-war-ticket-service/internal/app"
	"go-war-ticket-service/internal/platform/logger"
	"log"
)

func main() {
	// Load Configurations
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Setup Logger
	isDevelopment := cfg.ServerMode == "development"
	appLogger := logger.New(isDevelopment)
	defer appLogger.Sync() // Flush logs

	// Start Server
	srv := app.NewServer(cfg, appLogger)

	if err := srv.Start(); err != nil {
		appLogger.Fatal("server failed to start: ", err)
	}
}
