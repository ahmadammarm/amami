package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LimitBody caps the request body size; reads past the limit fail, which
// surfaces as a multipart/JSON parse error in the handler.
func LimitBody(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
