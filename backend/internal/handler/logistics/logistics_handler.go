package logistics

import (
	logisticsDto "github.com/ahmadammarm/amami/backend/internal/dto/logistics"
	logisticsSvc "github.com/ahmadammarm/amami/backend/internal/service/logistics"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LogisticsHandler struct {
	svc logisticsSvc.LogisticsService
}

func NewLogisticsHandler(svc logisticsSvc.LogisticsService) *LogisticsHandler {
	return &LogisticsHandler{svc}
}

func (h *LogisticsHandler) CreateAsset(c *gin.Context) {
	var req logisticsDto.CreateAssetRequest
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

func (h *LogisticsHandler) UpdateAsset(c *gin.Context) {
	id := c.Param("id")
	var req logisticsDto.CreateAssetRequest
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



func (h *LogisticsHandler) GetAllAssets(c *gin.Context) {
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

func (h *LogisticsHandler) DeleteAsset(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteAsset(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *LogisticsHandler) CreateLoan(c *gin.Context) {
	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset id"})
		return
	}
	var req logisticsDto.CreateAssetLoanRequest
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

func (h *LogisticsHandler) ReturnLoan(c *gin.Context) {
	loanID, err := uuid.Parse(c.Param("loanId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
		return
	}
	var req logisticsDto.ReturnAssetLoanRequest
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

func (h *LogisticsHandler) CreateAgenda(c *gin.Context) {
	var req logisticsDto.CreateAgendaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("userID")
	res, err := h.svc.CreateAgenda(req, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func (h *LogisticsHandler) GetAgendas(c *gin.Context) {
	res, err := h.svc.GetAgendas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *LogisticsHandler) UpdateAgenda(c *gin.Context) {
	id := c.Param("id")
	var req logisticsDto.CreateAgendaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.UpdateAgenda(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
