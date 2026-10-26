package inventory

import (
	inventoryDto "github.com/ahmadammarm/amami/backend/internal/dto/inventory"
	inventorySvc "github.com/ahmadammarm/amami/backend/internal/service/inventory"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	svc inventorySvc.InventoryService
}

func NewInventoryHandler(svc inventorySvc.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc}
}

func (h *InventoryHandler) CreateAsset(c *gin.Context) {
	var req inventoryDto.AssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.CreateAsset(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.CreatedResponse(c, "Asset created successfully", res)
}

func (h *InventoryHandler) UpdateAsset(c *gin.Context) {
	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req inventoryDto.AssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.UpdateAsset(id, req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Asset updated successfully", res)
}

func (h *InventoryHandler) GetAllAssets(c *gin.Context) {
	page, limit := utils.ParsePagination(c)

	assets, err := h.svc.GetAllAssets(page, limit)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Assets fetched successfully", assets)
}

func (h *InventoryHandler) DeleteAsset(c *gin.Context) {
	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err := h.svc.DeleteAsset(id); err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Asset deleted successfully", nil)
}

func (h *InventoryHandler) CreateLoan(c *gin.Context) {
	assetID, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req inventoryDto.CreateAssetLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.CreateLoan(assetID, req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.CreatedResponse(c, "Loan created successfully", res)
}

func (h *InventoryHandler) ReturnLoan(c *gin.Context) {
	loanID, err := utils.ParseUUIDParam(c, "loanId")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req inventoryDto.ReturnAssetLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.ReturnLoan(loanID, req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Loan returned successfully", res)
}
