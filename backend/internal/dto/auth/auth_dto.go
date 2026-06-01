package auth

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token                  string `json:"token"`
	RequiresPasswordChange bool   `json:"requires_password_change"`
}

type ChangePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
