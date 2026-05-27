package auth

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/dto/auth"
	"github.com/ahmadammarm/amami/backend/internal/dto/user"
	repo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(req auth.LoginRequest) (*auth.LoginResponse, error)
	GetMe(userID uuid.UUID) (*user.UserResponse, error)
}

type authService struct {
	userRepo repo.UserRepository
}

func NewAuthService(userRepo repo.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(req auth.LoginRequest) (*auth.LoginResponse, error) {
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

func (s *authService) GetMe(userID uuid.UUID) (*user.UserResponse, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user.UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		RoleName: u.Role.Name,
		Status:   u.Status,
	}, nil
}
