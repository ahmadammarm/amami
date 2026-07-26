package jamaah

import (
	"net/http"
	"strconv"

	"github.com/ahmadammarm/amami/backend/internal/dto/jamaah"
	jamaahSvc "github.com/ahmadammarm/amami/backend/internal/service/jamaah"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JamaahHandler struct {
	svc jamaahSvc.JamaahService
}

func NewJamaahHandler(svc jamaahSvc.JamaahService) *JamaahHandler {
	return &JamaahHandler{svc: svc}
}

func (h *JamaahHandler) CreateJamaah(c *gin.Context) {
	var req jamaah.CreateJamaahRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.CreateJamaah(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *JamaahHandler) UpdateJamaah(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var req jamaah.UpdateJamaahRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.UpdateJamaah(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *JamaahHandler) GetJamaahByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	res, err := h.svc.GetJamaahByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *JamaahHandler) GetJamaahs(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	query := c.Query("query")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	res, err := h.svc.GetJamaahs(page, limit, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
