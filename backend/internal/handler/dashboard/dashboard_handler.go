package dashboard

import (
	dashSvc "github.com/ahmadammarm/amami/backend/internal/service/dashboard"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
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
		utils.HandleError(c, apperror.Internal("failed to fetch dashboard metrics", err))
		return
	}
	utils.SuccessResponse(c, "Dashboard metrics fetched successfully", metrics)
}
