package zakat

import (
	"github.com/ahmadammarm/amami/backend/internal/dto/zakat"
	zakatSvc "github.com/ahmadammarm/amami/backend/internal/service/zakat"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ZakatHandler struct {
	svc zakatSvc.ZakatService
}

func NewZakatHandler(svc zakatSvc.ZakatService) *ZakatHandler {
	return &ZakatHandler{svc: svc}
}

func (h *ZakatHandler) CollectZakat(c *gin.Context) {
	userID, err := utils.CurrentUserID(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req zakat.CollectZakatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.CollectZakat(req, userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.CreatedResponse(c, "Zakat collected successfully", res)
}

func (h *ZakatHandler) DistributeZakat(c *gin.Context) {
	userID, err := utils.CurrentUserID(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req zakat.DistributeZakatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.DistributeZakat(req, userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.CreatedResponse(c, "Zakat distributed successfully", res)
}

func (h *ZakatHandler) GetDonations(c *gin.Context) {
	page, limit := utils.ParsePagination(c)

	res, err := h.svc.GetDonations(page, limit)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Donations fetched successfully", res)
}

func (h *ZakatHandler) GetDistributions(c *gin.Context) {
	page, limit := utils.ParsePagination(c)

	res, err := h.svc.GetDistributions(page, limit)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Distributions fetched successfully", res)
}
