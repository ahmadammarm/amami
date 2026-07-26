package inventory

import (
	inventoryDto "github.com/ahmadammarm/amami/backend/internal/dto/inventory"
	inventorySvc "github.com/ahmadammarm/amami/backend/internal/service/inventory"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	svc inventorySvc.InventoryService
}

func NewInventoryHandler(svc inventorySvc.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc}
}

func (h *InventoryHandler) CreateAsset(c *gin.Context) {
	var req inventoryDto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.CreateAsset(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func (h *InventoryHandler) UpdateAsset(c *gin.Context) {
	id := c.Param("id")
	var req inventoryDto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.UpdateAsset(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *InventoryHandler) GetAllAssets(c *gin.Context) {
	page := utils.ParseQueryInt(c, "page", 1)
	limit := utils.ParseQueryInt(c, "limit", 20)

	assets, total, err := h.svc.GetAllAssets(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch assets")
		return
	}
	
	utils.SuccessResponse(c, "Assets fetched successfully", gin.H{
		"items": assets,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *InventoryHandler) DeleteAsset(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteAsset(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *InventoryHandler) CreateLoan(c *gin.Context) {
	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset id"})
		return
	}
	var req inventoryDto.CreateAssetLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateLoan(assetID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "loan created"})
}

func (h *InventoryHandler) ReturnLoan(c *gin.Context) {
	loanID, err := uuid.Parse(c.Param("loanId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
		return
	}
	var req inventoryDto.ReturnAssetLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ReturnLoan(loanID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "loan returned"})
}
