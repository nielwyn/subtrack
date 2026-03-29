package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nielwyn/inventory-system/pkg/logger"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if rid, exists := c.Get(RequestIDKey); exists {
			fields = append(fields, zap.String("request_id", rid.(string)))
		}

		logger.Info("HTTP Request", fields...)

		for _, e := range c.Errors {
			logger.Error("Request error", zap.Error(e.Err))
		}
	}
}
