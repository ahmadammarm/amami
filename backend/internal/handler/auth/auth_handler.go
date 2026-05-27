package auth

import (
	"net/http"

	svc "github.com/ahmadammarm/amami/backend/internal/service/auth"
	"github.com/ahmadammarm/amami/backend/pkg/logger"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler interface {
	Login(c *gin.Context)
}

type authHandler struct {
	authService svc.AuthService
}

func NewAuthHandler(authService svc.AuthService) AuthHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) Login(c *gin.Context) {
	var req svc.LoginRequest
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
