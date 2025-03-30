package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}

	Server struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}

	Auth struct {
		JWTSecret   string
		JWTLifetime time.Duration
	}

	LogLevel string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	cfg := &Config{}

	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnv("DB_PORT", "5432")
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.DB.Name = getEnv("DB_NAME", "wishlist")
	cfg.DB.SSLMode = getEnv("DB_SSLMODE", "disable")

	cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	cfg.Server.ReadTimeout = 10 * time.Second
	cfg.Server.WriteTimeout = 10 * time.Second
	cfg.Server.IdleTimeout = 60 * time.Second

	cfg.Auth.JWTSecret = getEnv("JWT_SECRET", "default-secret")
	cfg.Auth.JWTLifetime = 24 * time.Hour

	cfg.LogLevel = getEnv("LOG_LEVEL", "info")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
