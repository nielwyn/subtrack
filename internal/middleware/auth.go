package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nielwyn/inventory-system/internal/service"
	apierrors "github.com/nielwyn/inventory-system/pkg/errors"
	"github.com/nielwyn/inventory-system/pkg/logger"
	"github.com/nielwyn/inventory-system/pkg/response"
	"go.uber.org/zap"
)

func Auth(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, apierrors.CodeUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, apierrors.CodeUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		token, err := authService.ValidateToken(parts[1])
		if err != nil {
			logger.Error("Token validation failed", zap.Error(err))
			response.Error(c, http.StatusUnauthorized, apierrors.CodeUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		userID, err := authService.GetUserFromToken(token)
		if err != nil {
			logger.Error("Failed to extract user from token", zap.Error(err))
			response.Error(c, http.StatusUnauthorized, apierrors.CodeUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
