package auth

import (
	"errors"

	repo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AuthService interface {
	Login(req LoginRequest) (*LoginResponse, error)
}

type authService struct {
	userRepo repo.UserRepository
}

func NewAuthService(userRepo repo.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(req LoginRequest) (*LoginResponse, error) {
	// Updated to specifically find by email
	user, err := s.userRepo.FindByUsernameOrEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.Status != "ACTIVE" {
		return nil, errors.New("account is suspended or inactive")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &LoginResponse{Token: token}, nil
}
