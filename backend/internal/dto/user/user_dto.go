package user

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
)

type InviteUserRequest struct {
	Username string `json:"username" binding:"required,max=50"`
	Email    string `json:"email" binding:"required,email,max=100"`
	FullName string `json:"full_name" binding:"required,max=100"`
	RoleID   uint   `json:"role_id" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateStatusRequest struct {
	Status domain.UserStatus `json:"status" binding:"required,oneof=ACTIVE SUSPENDED"`
}

type UserResponse struct {
	ID       uuid.UUID         `json:"id"`
	Username string            `json:"username"`
	Email    string            `json:"email"`
	FullName string            `json:"full_name"`
	RoleName string            `json:"role_name"`
	Status   domain.UserStatus `json:"status"`
}

// MeResponse is the authenticated user plus the permission codes of their role,
// which the frontend uses for route guards and UI visibility.
type MeResponse struct {
	UserResponse
	Permissions []domain.PermissionCode `json:"permissions"`
}

type InviteUserResponse struct {
	User              UserResponse `json:"user"`
	TemporaryPassword string       `json:"temporary_password"`
}

type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ToUserResponse maps a user (with Role and JamaahProfile preloaded) to its API shape.
func ToUserResponse(u *domain.User) UserResponse {
	res := UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		RoleName: u.Role.Name,
		Status:   u.Status,
	}
	if u.JamaahProfile != nil {
		res.FullName = u.JamaahProfile.FullName
	}
	return res
}
