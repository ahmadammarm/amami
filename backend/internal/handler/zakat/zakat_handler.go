package zakat

import (
	"net/http"
	"strconv"

	"github.com/ahmadammarm/amami/backend/internal/dto/zakat"
	zakatSvc "github.com/ahmadammarm/amami/backend/internal/service/zakat"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ZakatHandler struct {
	svc zakatSvc.ZakatService
}

func NewZakatHandler(svc zakatSvc.ZakatService) *ZakatHandler {
	return &ZakatHandler{svc: svc}
}

func (h *ZakatHandler) CollectZakat(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		parsed, _ := uuid.Parse(userIDStr.(string))
		userID = parsed
	}

	var req zakat.CollectZakatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.CollectZakat(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *ZakatHandler) DistributeZakat(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		parsed, _ := uuid.Parse(userIDStr.(string))
		userID = parsed
	}

	var req zakat.DistributeZakatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.DistributeZakat(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *ZakatHandler) GetDonations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	res, err := h.svc.GetDonations(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ZakatHandler) GetDistributions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	res, err := h.svc.GetDistributions(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
