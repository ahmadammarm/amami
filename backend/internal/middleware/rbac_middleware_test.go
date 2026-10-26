package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type fakeFinder struct {
	perms map[uint][]domain.Permission
	err   error
}

func (f fakeFinder) FindPermissionsByRoleID(roleID uint) ([]domain.Permission, error) {
	return f.perms[roleID], f.err
}

func TestRequirePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	finder := fakeFinder{perms: map[uint][]domain.Permission{
		1: {{Code: string(domain.PermLedgerRead)}, {Code: string(domain.PermLedgerWrite)}},
		2: {{Code: string(domain.PermLedgerRead)}},
	}}

	tests := []struct {
		name       string
		finder     PermissionFinder
		roleID     any
		wantStatus int
	}{
		{name: "granted", finder: finder, roleID: uint(1), wantStatus: http.StatusOK},
		{name: "missing permission", finder: finder, roleID: uint(2), wantStatus: http.StatusForbidden},
		{name: "unknown role", finder: finder, roleID: uint(99), wantStatus: http.StatusForbidden},
		{name: "no role in context", finder: finder, roleID: nil, wantStatus: http.StatusUnauthorized},
		{name: "lookup failure", finder: fakeFinder{err: errors.New("db down")}, roleID: uint(1), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/",
				func(c *gin.Context) {
					if tt.roleID != nil {
						c.Set(utils.ContextKeyRoleID, tt.roleID)
					}
				},
				RequirePermission(tt.finder, domain.PermLedgerWrite),
				func(c *gin.Context) { c.Status(http.StatusOK) },
			)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
