package user

import (
	"github.com/ahmadammarm/amami/backend/internal/dto/user"
	userSvc "github.com/ahmadammarm/amami/backend/internal/service/user"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserMgmtHandler interface {
	InviteUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	GetRoles(c *gin.Context)
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

	res, err := h.svc.InviteUser(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.CreatedResponse(c, "User invited successfully", res)
}

func (h *userMgmtHandler) GetAllUsers(c *gin.Context) {
	page, limit := utils.ParsePagination(c)

	users, err := h.svc.GetAllUsers(page, limit)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, "Users fetched successfully", users)
}

func (h *userMgmtHandler) GetRoles(c *gin.Context) {
	roles, err := h.svc.GetRoles()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, "Roles fetched successfully", roles)
}

func (h *userMgmtHandler) UpdateUserStatus(c *gin.Context) {
	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req user.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.svc.UpdateUserStatus(id, req.Status); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, "User status updated successfully", nil)
}

func (h *userMgmtHandler) DeleteUser(c *gin.Context) {
	actorID, err := utils.CurrentUserID(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err := h.svc.DeleteUser(actorID, id); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, "User deleted successfully", nil)
}
