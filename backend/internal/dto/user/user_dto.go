package user

import "github.com/google/uuid"

type InviteUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string    `json:"email" binding:"required,email"`
	FullName string    `json:"full_name" binding:"required"`
	RoleID   uint      `json:"role_id" binding:"required"`
	Password string    `json:"password" binding:"required,min=8"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=ACTIVE SUSPENDED"`
}


type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	RoleName string    `json:"role_name"`
	Status   string    `json:"status"`
}
