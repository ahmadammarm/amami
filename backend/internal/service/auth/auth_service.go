package auth

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/dto/auth"
	repo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
)

type AuthService interface {
	Login(req auth.LoginRequest) (*auth.LoginResponse, error)
}

type authService struct {
	userRepo repo.UserRepository
}

func NewAuthService(userRepo repo.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(req auth.LoginRequest) (*auth.LoginResponse, error) {
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

	return &auth.LoginResponse{Token: token}, nil
}
