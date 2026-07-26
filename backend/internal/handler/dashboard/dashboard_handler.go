package dashboard

import (
	"net/http"

	dashSvc "github.com/ahmadammarm/amami/backend/internal/service/dashboard"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	svc dashSvc.DashboardService
}

func NewDashboardHandler(svc dashSvc.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetMetrics(c *gin.Context) {
	metrics, err := h.svc.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}
