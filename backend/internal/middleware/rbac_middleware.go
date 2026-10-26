package middleware

import (
	"net/http"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PermissionFinder is the subset of auth.RoleRepository the RBAC middleware needs.
type PermissionFinder interface {
	FindPermissionsByRoleID(roleID uint) ([]domain.Permission, error)
}

func RequirePermission(finder PermissionFinder, code domain.PermissionCode) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, err := utils.CurrentRoleID(c)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Role ID not found in context")
			c.Abort()
			return
		}

		permissions, err := finder.FindPermissionsByRoleID(roleID)
		if err != nil {
			zap.L().Error("failed to fetch permissions", zap.Uint("role_id", roleID), zap.Error(err))
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch permissions")
			c.Abort()
			return
		}

		for _, p := range permissions {
			if p.Code == string(code) {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusForbidden, "Permission denied")
		c.Abort()
	}
}
