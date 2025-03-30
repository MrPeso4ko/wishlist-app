package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"wishlist-app/internal/config"
	"wishlist-app/internal/models"
	"wishlist-app/pkg/logger"
)

func NewDB(cfg *config.Config, log logger.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.Port,
		cfg.DB.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn))

	logger, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if cfg.LogLevel == "debug" {
		logger.Info("Database debug logging enabled")
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Wish{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	logger.Info("Database connection established and migrations applied")
	return db, nil
}
