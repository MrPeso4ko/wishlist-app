package handler

import (
	"net/http"

	"wishlist-app/internal/config"
	"wishlist-app/internal/service"
	"wishlist-app/pkg/logger"
	"wishlist-app/pkg/metrics"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      logger.Logger
	cfg         *config.Config
}

func NewAuthHandler(cfg *config.Config, logger logger.Logger, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
		logger:      logger,
	}
}

type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// @BasePath /api

// Register godoc
// @Summary Register a new user
// @Description Register a new user with login and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 201 "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.RecordAuthRequest("register", "failure")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.Register(req.Login, req.Password); err != nil {
		metrics.RecordAuthRequest("register", "failure")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metrics.RecordAuthRequest("register", "success")
	c.Status(http.StatusCreated)
}

// Login godoc
// @Summary Login a user
// @Description Login a user with login and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} LoginResponse "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.RecordAuthRequest("login", "failure")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(req.Login, req.Password)
	if err != nil {
		metrics.RecordAuthRequest("login", "failure")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	metrics.RecordAuthRequest("login", "success")
	c.JSON(http.StatusOK, LoginResponse{Token: token})
}
