package jamaah

import (
	"github.com/ahmadammarm/amami/backend/internal/dto/jamaah"
	jamaahSvc "github.com/ahmadammarm/amami/backend/internal/service/jamaah"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type JamaahHandler struct {
	svc jamaahSvc.JamaahService
}

func NewJamaahHandler(svc jamaahSvc.JamaahService) *JamaahHandler {
	return &JamaahHandler{svc: svc}
}

func (h *JamaahHandler) CreateJamaah(c *gin.Context) {
	var req jamaah.JamaahRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.CreateJamaah(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.CreatedResponse(c, "Jamaah created successfully", res)
}

func (h *JamaahHandler) UpdateJamaah(c *gin.Context) {
	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req jamaah.JamaahRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.UpdateJamaah(id, req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Jamaah updated successfully", res)
}

func (h *JamaahHandler) GetJamaahByID(c *gin.Context) {
	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	res, err := h.svc.GetJamaahByID(id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Jamaah fetched successfully", res)
}

func (h *JamaahHandler) GetJamaahs(c *gin.Context) {
	page, limit := utils.ParsePagination(c)

	var query jamaah.ListJamaahQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.GetJamaahs(page, limit, query)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Jamaah fetched successfully", res)
}

func (h *JamaahHandler) GetSurvey(c *gin.Context) {
	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	res, err := h.svc.GetSurvey(id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Survey fetched successfully", res)
}

func (h *JamaahHandler) SubmitSurvey(c *gin.Context) {
	surveyorID, err := utils.CurrentUserID(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req jamaah.MustahikSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.SubmitSurvey(id, req, surveyorID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Survey saved successfully", res)
}
