package middleware

import (
	"github.com/ahmadammarm/amami/backend/internal/repository/auth"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func RequirePermission(roleRepo auth.RoleRepository, code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			utils.ErrorResponse(c, 401, "Role ID not found in context")
			c.Abort()
			return
		}

		permissions, err := roleRepo.FindPermissionsByRoleID(roleID.(uint))
		if err != nil {
			utils.ErrorResponse(c, 500, "Failed to fetch permissions")
			c.Abort()
			return
		}

		hasPermission := false
		for _, p := range permissions {
			if p.Code == code {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			utils.ErrorResponse(c, 403, "Permission denied")
			c.Abort()
			return
		}

		c.Next()
	}
}
