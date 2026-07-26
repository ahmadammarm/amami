package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
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
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ParseQueryInt parses a query parameter to an integer with a fallback default value
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

func ErrorResponse(c *gin.Context, code int, message string) {
	JSONResponse(c, code, false, message, nil, nil)
}

func ValidationErrorResponse(c *gin.Context, err error) {
	errs := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
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
	}
	return fmt.Sprintf("%s is invalid", e.Field())
}
