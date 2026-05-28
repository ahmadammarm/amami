package settings

import (
	"net/http"
	"strconv"

	"github.com/ahmadammarm/amami/backend/internal/dto/settings"
	settingsSvc "github.com/ahmadammarm/amami/backend/internal/service/settings"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SettingsHandler interface {
	GetMosqueProfile(c *gin.Context)
	UpdateMosqueProfile(c *gin.Context)
	UpdateSMTPConfig(c *gin.Context)
	GetSMTPConfig(c *gin.Context)
	TestSMTPConnection(c *gin.Context)
	GetSystemHealth(c *gin.Context)
	GetAuditLogs(c *gin.Context)
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

func (h *settingsHandler) GetSMTPConfig(c *gin.Context) {
	config, err := h.svc.GetSMTPConfig()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch SMTP configuration")
		return
	}

	utils.SuccessResponse(c, "SMTP configuration fetched successfully", config)
}

func (h *settingsHandler) TestSMTPConnection(c *gin.Context) {
	var req settings.TestSMTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.svc.TestSMTPConnection(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, "SMTP connection successful", nil)
}

func (h *settingsHandler) GetSystemHealth(c *gin.Context) {
	health, err := h.svc.GetSystemHealth()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch system health")
		return
	}
	utils.SuccessResponse(c, "System health fetched successfully", health)
}

func (h *settingsHandler) GetAuditLogs(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	logs, err := h.svc.GetAuditLogs(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch audit logs")
		return
	}
	utils.SuccessResponse(c, "Audit logs fetched successfully", logs)
}

func (h *settingsHandler) UpdateMosqueProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req settings.MosqueProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.svc.UpdateMosqueProfile(userID, req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update mosque profile")
		return
	}

	utils.SuccessResponse(c, "Mosque profile updated successfully", nil)
}

func (h *settingsHandler) UpdateSMTPConfig(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req settings.SMTPConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.svc.UpdateSMTPConfig(userID, req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update SMTP config")
		return
	}

	utils.SuccessResponse(c, "SMTP configuration updated successfully", nil)
}
