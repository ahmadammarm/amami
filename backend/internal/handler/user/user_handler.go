package user

import (
	"net/http"

	"github.com/ahmadammarm/amami/backend/internal/dto/user"
	userSvc "github.com/ahmadammarm/amami/backend/internal/service/user"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserMgmtHandler interface {
	InviteUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	UpdateUserStatus(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userMgmtHandler struct {
	svc userSvc.UserMgmtService
}

func NewUserMgmtHandler(svc userSvc.UserMgmtService) UserMgmtHandler {
	return &userMgmtHandler{svc: svc}
}

func (h *userMgmtHandler) InviteUser(c *gin.Context) {
	var req user.InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userRes, tempPass, err := h.svc.InviteUser(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// For now, we return the temp password in the response until SMTP is ready
	utils.SuccessResponse(c, "User invited successfully", gin.H{
		"user":              userRes,
		"temporary_password": tempPass,
	})
}

func (h *userMgmtHandler) GetAllUsers(c *gin.Context) {
	page := utils.ParseQueryInt(c, "page", 1)
	limit := utils.ParseQueryInt(c, "limit", 20)

	users, total, err := h.svc.GetAllUsers(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	utils.SuccessResponse(c, "Users fetched successfully", gin.H{
		"items": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *userMgmtHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id format")
		return
	}

	var req user.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.svc.UpdateUserStatus(id, req.Status); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, "User status updated successfully", nil)
}

func (h *userMgmtHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id format")
		return
	}

	if err := h.svc.DeleteUser(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, "User deleted successfully", nil)
}
