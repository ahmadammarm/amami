package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestHandleError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		err         error
		wantStatus  int
		wantMessage string
	}{
		{"validation", apperror.Validation("bad input"), http.StatusBadRequest, "bad input"},
		{"not found", apperror.NotFound("missing"), http.StatusNotFound, "missing"},
		{"conflict", apperror.Conflict("taken"), http.StatusConflict, "taken"},
		{"unauthorized", apperror.Unauthorized("who are you"), http.StatusUnauthorized, "who are you"},
		{"forbidden", apperror.Forbidden("nope"), http.StatusForbidden, "nope"},
		{"wrapped app error", fmt.Errorf("context: %w", apperror.Conflict("taken")), http.StatusConflict, "taken"},
		{"internal hides cause", apperror.Internal("failed to save", errors.New("pq: secret detail")), http.StatusInternalServerError, "failed to save"},
		{"gorm not found", gorm.ErrRecordNotFound, http.StatusNotFound, "resource not found"},
		{"unknown error is masked", errors.New("raw driver error"), http.StatusInternalServerError, "internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			HandleError(c, tt.err)

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			var body Response
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("invalid JSON body: %v", err)
			}
			if body.Success || body.Message != tt.wantMessage {
				t.Fatalf("body = %+v, want message %q", body, tt.wantMessage)
			}
		})
	}
}

func TestNewPage(t *testing.T) {
	page := NewPage[int](nil, 21, 2, 10)
	if page.TotalPages != 3 {
		t.Fatalf("TotalPages = %d, want 3", page.TotalPages)
	}
	if page.Items == nil {
		t.Fatalf("Items must be an empty slice, not nil, so it encodes as []")
	}
}
