package utils

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	ContextKeyUserID     = "user_id"
	ContextKeyRoleID     = "role_id"
	ContextKeyTokenScope = "token_scope"

	defaultPageLimit = 20
	maxPageLimit     = 100
)

// UseJSONFieldNames makes validation errors report json tag names (e.g. "fund_id")
// instead of Go struct field names, so the frontend can map errors to form fields.
func UseJSONFieldNames() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		for _, tag := range []string{"json", "form"} {
			name := strings.SplitN(field.Tag.Get(tag), ",", 2)[0]
			if name != "" && name != "-" {
				return name
			}
		}
		return field.Name
	})
}

func CurrentUserID(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(ContextKeyUserID)
	if !ok {
		return uuid.Nil, apperror.Unauthorized("user is not authenticated")
	}
	userID, ok := id.(uuid.UUID)
	if !ok {
		return uuid.Nil, apperror.Internal("invalid user id in context", nil)
	}
	return userID, nil
}

func CurrentRoleID(c *gin.Context) (uint, error) {
	id, ok := c.Get(ContextKeyRoleID)
	if !ok {
		return 0, apperror.Unauthorized("user is not authenticated")
	}
	roleID, ok := id.(uint)
	if !ok {
		return 0, apperror.Internal("invalid role id in context", nil)
	}
	return roleID, nil
}

// ParseQueryInt parses a query parameter to an integer with a fallback default value.
func ParseQueryInt(c *gin.Context, key string, defaultValue int) int {
	valStr := c.Query(key)
	if valStr == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

// ParsePagination reads ?page and ?limit, clamping them to sane bounds.
func ParsePagination(c *gin.Context) (page, limit int) {
	page = ParseQueryInt(c, "page", 1)
	if page < 1 {
		page = 1
	}
	limit = ParseQueryInt(c, "limit", defaultPageLimit)
	if limit < 1 || limit > maxPageLimit {
		limit = defaultPageLimit
	}
	return page, limit
}

func ParseUUIDParam(c *gin.Context, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		return uuid.Nil, apperror.Validation("invalid " + name)
	}
	return id, nil
}

func ParseUintParam(c *gin.Context, name string) (uint, error) {
	id, err := strconv.ParseUint(c.Param(name), 10, 32)
	if err != nil {
		return 0, apperror.Validation("invalid " + name)
	}
	return uint(id), nil
}
