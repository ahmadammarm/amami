package qurban

import (
	qurbanDto "github.com/ahmadammarm/amami/backend/internal/dto/qurban"
	qurbanSvc "github.com/ahmadammarm/amami/backend/internal/service/qurban"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type QurbanHandler struct {
	svc qurbanSvc.QurbanService
}

func NewQurbanHandler(svc qurbanSvc.QurbanService) *QurbanHandler {
	return &QurbanHandler{svc}
}

func (h *QurbanHandler) CreatePackage(c *gin.Context) {
	var req qurbanDto.PackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.CreatePackage(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.CreatedResponse(c, "Package created successfully", res)
}

func (h *QurbanHandler) GetAllPackages(c *gin.Context) {
	page, limit := utils.ParsePagination(c)

	pkgs, err := h.svc.GetAllPackages(page, limit)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Packages fetched successfully", pkgs)
}

func (h *QurbanHandler) UpdatePackage(c *gin.Context) {
	id, err := utils.ParseUintParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req qurbanDto.PackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.UpdatePackage(id, req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Package updated successfully", res)
}

func (h *QurbanHandler) DeletePackage(c *gin.Context) {
	id, err := utils.ParseUintParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err := h.svc.DeletePackage(id); err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Package deleted successfully", nil)
}

func (h *QurbanHandler) CreateBooking(c *gin.Context) {
	var req qurbanDto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.CreateBooking(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.CreatedResponse(c, "Booking created successfully", res)
}

func (h *QurbanHandler) GetBookings(c *gin.Context) {
	res, err := h.svc.GetBookings()
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Bookings fetched successfully", res)
}

func (h *QurbanHandler) CreateAnimal(c *gin.Context) {
	var req qurbanDto.AnimalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.CreateAnimal(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.CreatedResponse(c, "Animal created successfully", res)
}

func (h *QurbanHandler) GetAnimals(c *gin.Context) {
	res, err := h.svc.GetAnimals()
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Animals fetched successfully", res)
}

func (h *QurbanHandler) UpdateAnimal(c *gin.Context) {
	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req qurbanDto.AnimalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.UpdateAnimal(id, req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Animal updated successfully", res)
}

func (h *QurbanHandler) DeleteAnimal(c *gin.Context) {
	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err := h.svc.DeleteAnimal(id); err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Animal deleted successfully", nil)
}
