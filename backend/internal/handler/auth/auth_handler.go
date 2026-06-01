package auth

import (
	"net/http"

	"github.com/ahmadammarm/amami/backend/internal/dto/auth"
	svc "github.com/ahmadammarm/amami/backend/internal/service/auth"
	"github.com/ahmadammarm/amami/backend/pkg/logger"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthHandler interface {
	Login(c *gin.Context)
	GetMe(c *gin.Context)
	ChangePassword(c *gin.Context)
}

type authHandler struct {
	authService svc.AuthService
}

func NewAuthHandler(authService svc.AuthService) AuthHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Login validation failed", zap.Error(err))
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		logger.Info("Login attempt failed", 
			zap.String("email", req.Email), 
			zap.String("reason", err.Error()),
		)
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	logger.Info("Login successful", zap.String("email", req.Email))
	utils.SuccessResponse(c, "Login successful", res)
}

func (h *authHandler) ChangePassword(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	var req auth.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.authService.ChangePassword(userID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, "Password changed successfully", res)
}

func (h *authHandler) GetMe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "invalid user id type in context")
		return
	}

	user, err := h.authService.GetMe(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, "User info fetched successfully", user)
}
