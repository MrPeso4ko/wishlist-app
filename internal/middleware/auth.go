package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"wishlist-app/internal/config"
	"wishlist-app/internal/service"
	"wishlist-app/pkg/logger"
)

func Auth(cfg *config.Config, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bearer token required"})
			return
		}

		authService := service.NewAuthService(nil, cfg)
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			logger.Warnf("Invalid token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
