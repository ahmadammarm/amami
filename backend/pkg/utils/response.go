package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Page is the standard envelope for paginated list data.
type Page[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

func NewPage[T any](items []T, total int64, page, limit int) Page[T] {
	if items == nil {
		items = []T{}
	}
	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}
	return Page[T]{Items: items, Total: total, Page: page, Limit: limit, TotalPages: totalPages}
}

func JSONResponse(c *gin.Context, code int, success bool, message string, data interface{}, errs interface{}) {
	c.JSON(code, Response{
		Success: success,
		Message: message,
		Data:    data,
		Errors:  errs,
	})
}

func SuccessResponse(c *gin.Context, message string, data interface{}) {
	JSONResponse(c, http.StatusOK, true, message, data, nil)
}

func CreatedResponse(c *gin.Context, message string, data interface{}) {
	JSONResponse(c, http.StatusCreated, true, message, data, nil)
}

func ErrorResponse(c *gin.Context, code int, message string) {
	JSONResponse(c, code, false, message, nil, nil)
}

// HandleError writes the response for an error returned by a service.
// Unknown errors are logged and reported as a generic 500 so internals never leak.
func HandleError(c *gin.Context, err error) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		status := statusForKind(appErr.Kind)
		if status == http.StatusInternalServerError {
			zap.L().Error(appErr.Msg, zap.Error(appErr.Err), zap.String("path", c.FullPath()))
		}
		ErrorResponse(c, status, appErr.Msg)
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ErrorResponse(c, http.StatusNotFound, "resource not found")
		return
	}

	zap.L().Error("unhandled error", zap.Error(err), zap.String("path", c.FullPath()))
	ErrorResponse(c, http.StatusInternalServerError, "internal server error")
}

func statusForKind(kind apperror.Kind) int {
	switch kind {
	case apperror.KindValidation:
		return http.StatusBadRequest
	case apperror.KindNotFound:
		return http.StatusNotFound
	case apperror.KindConflict:
		return http.StatusConflict
	case apperror.KindUnauthorized:
		return http.StatusUnauthorized
	case apperror.KindForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func ValidationErrorResponse(c *gin.Context, err error) {
	errs := make(map[string]string)

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			errs[e.Field()] = formatValidationError(e)
		}
	} else {
		errs["error"] = err.Error()
	}

	JSONResponse(c, http.StatusBadRequest, false, "Validation failed", nil, errs)
}

func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	case "gt", "gte", "lt", "lte":
		return fmt.Sprintf("%s must be %s %s", e.Field(), e.Tag(), e.Param())
	}
	return fmt.Sprintf("%s is invalid", e.Field())
}
