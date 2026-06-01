package middleware

import (
	"strings"

	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, 401, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, 401, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			utils.ErrorResponse(c, 401, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role_id", claims.RoleID)
		c.Set("token_scope", claims.Scope)

		// By default, most routes require full_access.
		// If a token is restricted to password_reset, it should only hit the change-password route.
		// We'll handle specific overrides where needed or check it here.
		c.Next()
	}
}

func RequireScope(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scope, exists := c.Get("token_scope")
		if !exists || scope != requiredScope {
			utils.ErrorResponse(c, 403, "Insufficient token scope")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RejectScope(rejectedScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scope, exists := c.Get("token_scope")
		if exists && scope == rejectedScope {
			utils.ErrorResponse(c, 403, "This action is restricted until password is changed")
			c.Abort()
			return
		}
		c.Next()
	}
}
