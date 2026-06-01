package auth

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/dto/auth"
	"github.com/ahmadammarm/amami/backend/internal/dto/user"
	repo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(req auth.LoginRequest) (*auth.LoginResponse, error)
	GetMe(userID uuid.UUID) (*user.UserResponse, error)
	ChangePassword(userID uuid.UUID, req auth.ChangePasswordRequest) (*auth.LoginResponse, error)
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

	if user.Status != "ACTIVE" && user.Status != "PENDING_PASSWORD_CHANGE" {
		return nil, errors.New("account is suspended or inactive")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	scope := utils.ScopeFullAccess
	requiresChange := false
	if user.Status == "PENDING_PASSWORD_CHANGE" {
		scope = utils.ScopePasswordReset
		requiresChange = true
	}

	token, err := utils.GenerateToken(user.ID, user.RoleID, scope)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Record Audit Log
	s.userRepo.CreateAuditLog(domain.AuditLog{
		UserID: user.ID,
		Action: "LOGIN",
		Entity: "User Session",
	})

	return &auth.LoginResponse{
		Token:                  token,
		RequiresPasswordChange: requiresChange,
	}, nil
}

func (s *authService) ChangePassword(userID uuid.UUID, req auth.ChangePasswordRequest) (*auth.LoginResponse, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	u.PasswordHash = hashedPassword
	u.Status = "ACTIVE"

	if err := s.userRepo.Update(u); err != nil {
		return nil, errors.New("failed to update user")
	}

	token, err := utils.GenerateToken(u.ID, u.RoleID, utils.ScopeFullAccess)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	s.userRepo.CreateAuditLog(domain.AuditLog{
		UserID: u.ID,
		Action: "PASSWORD_CHANGE",
		Entity: "User Account",
	})

	return &auth.LoginResponse{
		Token:                  token,
		RequiresPasswordChange: false,
	}, nil
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
