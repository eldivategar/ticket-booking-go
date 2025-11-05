package database

import (
	"fmt"
	"go-service-boilerplate/configs"
	"go-service-boilerplate/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate (Hanya untuk development, idealnya gunakan file migrasi)
	if err = db.AutoMigrate(&domain.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
