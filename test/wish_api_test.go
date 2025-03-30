package test

import (
	"bytes"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"wishlist-app/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"wishlist-app/internal/config"
	"wishlist-app/internal/handler"
	"wishlist-app/internal/models"
	"wishlist-app/internal/service"
	"wishlist-app/pkg/logger"
)

var (
	mockWishRepo = new(MockWishRepository)
	mockUserRepo = new(MockUserRepository)
)

func setupWishRouter() *gin.Engine {
	cfg := &config.Config{}
	cfg.Auth.JWTSecret = "test-secret"
	cfg.Auth.JWTLifetime = 24 * time.Hour
	log, _ := logger.New("test")

	authService := service.NewAuthService(mockUserRepo, cfg)
	wishService := service.NewWishService(mockWishRepo, mockUserRepo)

	router := gin.New()
	authHandler := handler.NewAuthHandler(cfg, log, authService)
	wishHandler := handler.NewWishHandler(cfg, log, wishService)

	api := router.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.GET("/wishes/:username", wishHandler.GetByUsername)

		auth := api.Group("")
		auth.Use(middleware.Auth(cfg, log))
		{
			auth.POST("/wishes", wishHandler.Create)
			auth.PUT("/wishes/:id", wishHandler.Update)
			auth.DELETE("/wishes/:id", wishHandler.Delete)
			auth.GET("/wishes", wishHandler.GetByUserID)
		}
	}

	return router
}

func TestCreateWish(t *testing.T) {
	router := setupWishRouter()

	token := getTestToken(t, router)

	wish := models.Wish{
		Title: "Test Wish",
	}
	body, _ := json.Marshal(wish)

	w := httptest.NewRecorder()
	mockWishRepo.On("Create", &wish).Return(nil)

	req, _ := http.NewRequest("POST", "/api/wishes", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetPublicWishes(t *testing.T) {
	router := setupWishRouter()

	mockWishRepo.On("GetByUsername", "testuser").Return([]models.Wish{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/wishes/testuser", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func getTestToken(t *testing.T, router *gin.Engine) string {
	creds := map[string]string{
		"login":    "testuser",
		"password": "testpass",
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	mockUserRepo.On("FindByLogin", "testuser").Return(&models.User{Login: "testuser", PasswordHash: string(hashedPassword)}, nil)

	body, _ := json.Marshal(creds)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	return response["token"]
}
