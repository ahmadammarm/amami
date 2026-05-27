package settings

import (
	"net/http"

	"github.com/ahmadammarm/amami/backend/internal/dto/settings"
	settingsSvc "github.com/ahmadammarm/amami/backend/internal/service/settings"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type SettingsHandler interface {
	GetMosqueProfile(c *gin.Context)
	UpdateMosqueProfile(c *gin.Context)
	UpdateSMTPConfig(c *gin.Context)
}

type settingsHandler struct {
	svc settingsSvc.SettingsService
}

func NewSettingsHandler(svc settingsSvc.SettingsService) SettingsHandler {
	return &settingsHandler{svc: svc}
}

func (h *settingsHandler) GetMosqueProfile(c *gin.Context) {
	profile, err := h.svc.GetMosqueProfile()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch mosque profile")
		return
	}

	utils.SuccessResponse(c, "Profile fetched successfully", profile)
}

func (h *settingsHandler) UpdateMosqueProfile(c *gin.Context) {
	var req settings.MosqueProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.svc.UpdateMosqueProfile(req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update mosque profile")
		return
	}

	utils.SuccessResponse(c, "Profile updated successfully", nil)
}

func (h *settingsHandler) UpdateSMTPConfig(c *gin.Context) {
	var req settings.SMTPConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.svc.UpdateSMTPConfig(req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update SMTP config")
		return
	}

	utils.SuccessResponse(c, "SMTP configuration updated successfully", nil)
}
