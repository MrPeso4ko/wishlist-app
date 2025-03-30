package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"wishlist-app/pkg/logger"
)

func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				log.Error(e)
			}
		} else {
			log.Infof("[%d] %s %s %s %s %s",
				c.Writer.Status(),
				c.Request.Method,
				path,
				query,
				c.ClientIP(),
				latency,
			)
		}
	}
}
