package qurban

import (
	qurbanDto "github.com/ahmadammarm/amami/backend/internal/dto/qurban"
	qurbanSvc "github.com/ahmadammarm/amami/backend/internal/service/qurban"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QurbanHandler struct {
	svc qurbanSvc.QurbanService
}

func NewQurbanHandler(svc qurbanSvc.QurbanService) *QurbanHandler {
	return &QurbanHandler{svc}
}

func (h *QurbanHandler) CreatePackage(c *gin.Context) {
	var req qurbanDto.CreatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.CreatePackage(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func (h *QurbanHandler) GetAllPackages(c *gin.Context) {
	page := utils.ParseQueryInt(c, "page", 1)
	limit := utils.ParseQueryInt(c, "limit", 20)

	pkgs, total, err := h.svc.GetAllPackages(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch packages")
		return
	}
	
	utils.SuccessResponse(c, "Packages fetched successfully", gin.H{
		"items": pkgs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *QurbanHandler) UpdatePackage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid package id"})
		return
	}

	var req qurbanDto.CreatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.UpdatePackage(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *QurbanHandler) DeletePackage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid package id"})
		return
	}

	if err := h.svc.DeletePackage(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *QurbanHandler) CreateBooking(c *gin.Context) {
	var req qurbanDto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.CreateBooking(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func (h *QurbanHandler) GetBookings(c *gin.Context) {
	res, err := h.svc.GetBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *QurbanHandler) CreateAnimal(c *gin.Context) {
	var req qurbanDto.CreateAnimalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.CreateAnimal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func (h *QurbanHandler) GetAnimals(c *gin.Context) {
	res, err := h.svc.GetAnimals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *QurbanHandler) UpdateAnimal(c *gin.Context) {
	id := c.Param("id")
	var req qurbanDto.CreateAnimalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.UpdateAnimal(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *QurbanHandler) DeleteAnimal(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteAnimal(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
