// @title Wishlist API
// @version 1.0
// @description API server for managing wishes.
// @termsOfService http://swagger.io/terms/
// @contact.name Petr Salnikov
// @contact.email pdsalnikov@edu.hse.ru
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"wishlist-app/internal/config"

	"wishlist-app/internal/server"

	"wishlist-app/pkg/logger"
	"wishlist-app/pkg/metrics"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	metrics.Init()

	srv := server.New(cfg, logger)
	go func() {
		if err := srv.Run(); err != nil {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("Shutting down server...")
	if err := srv.Shutdown(); err != nil {
		logger.Errorf("Error during server shutdown: %v", err)
	}
}
