package server

import (
	"net/http"
	"wishlist-app/internal/repository"
	"wishlist-app/internal/service"

	_ "wishlist-app/docs"
	"wishlist-app/internal/config"
	"wishlist-app/internal/handler"
	"wishlist-app/internal/middleware"
	"wishlist-app/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
}

func New(cfg *config.Config, logger logger.Logger) *Server {
	router := gin.New()

	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.Metrics())

	db, err := repository.NewDB(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	wishRepo := repository.NewWishRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	wishService := service.NewWishService(wishRepo, userRepo)

	api := router.Group("/api")
	{
		authHandler := handler.NewAuthHandler(cfg, logger, authService)
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		wishHandler := handler.NewWishHandler(cfg, logger, wishService)
		api.GET("/wishes/:username", wishHandler.GetByUsername)

		auth := api.Group("")
		auth.Use(middleware.Auth(cfg, logger))
		{
			auth.POST("/wishes", wishHandler.Create)
			auth.PUT("/wishes/:id", wishHandler.Update)
			auth.DELETE("/wishes/:id", wishHandler.Delete)
			auth.GET("/wishes", wishHandler.GetByUserID)
		}
	}

	// Metrics endpoint
	// @Summary Prometheus metrics endpoint
	// @Description Get Prometheus metrics
	// @Tags metrics
	// @Produce text/plain
	// @Success 200 {string} string "OK"
	// @Router /metrics [get]
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check
	// @Summary Health check endpoint
	// @Description Check if the server is running
	// @Tags health
	// @Produce json
	// @Success 200 {object} map[string]string "OK"
	// @Router /health [get]
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.Server.Port,
			Handler:      router,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		logger: logger,
	}
}

func (s *Server) Run() error {
	s.logger.Infof("Server is running on port %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	s.logger.Info("Shutting down server...")
	return s.httpServer.Close()
}
