package user

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/dto/user"
	authRepo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	userRepo "github.com/ahmadammarm/amami/backend/internal/repository/user"
	"github.com/ahmadammarm/amami/backend/pkg/logger"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserMgmtService interface {
	InviteUser(req user.InviteUserRequest) (*user.UserResponse, string, error)
	GetAllUsers() ([]user.UserResponse, error)
	UpdateUserStatus(userID uuid.UUID, status string) error
}

type userMgmtService struct {
	userRepo userRepo.UserMgmtRepository
	roleRepo authRepo.RoleRepository
}

func NewUserMgmtService(userRepo userRepo.UserMgmtRepository, roleRepo authRepo.RoleRepository) UserMgmtService {
	return &userMgmtService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *userMgmtService) InviteUser(req user.InviteUserRequest) (*user.UserResponse, string, error) {
	// 1. Validate Role
	role, err := s.roleRepo.FindByID(req.RoleID)
	if err != nil {
		return nil, "", errors.New("invalid role id")
	}

	// 2. Generate random password
	tempPassword, err := generateRandomPassword(12)
	if err != nil {
		return nil, "", errors.New("failed to generate temporary password")
	}

	hashedPassword, err := utils.HashPassword(tempPassword)
	if err != nil {
		return nil, "", errors.New("failed to hash password")
	}

	// 3. Construct models
	userModel := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		RoleID:       req.RoleID,
		Status:       "ACTIVE",
	}

	jamaah := &domain.Jamaah{
		FullName: req.FullName,
	}

	// 4. Save to DB
	if err := s.userRepo.CreateUserWithProfile(userModel, jamaah); err != nil {
		logger.Error("failed to create user with profile", zap.Error(err))
		return nil, "", errors.New("failed to create user (username or email might already exist)")
	}

	logger.Info("User invited successfully", 
		zap.String("username", userModel.Username), 
		zap.String("temp_password", tempPassword), // Log it until SMTP is implemented
	)

	return &user.UserResponse{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		FullName: jamaah.FullName,
		RoleName: role.Name,
		Status:   userModel.Status,
	}, tempPassword, nil
}

func (s *userMgmtService) GetAllUsers() ([]user.UserResponse, error) {
	users, err := s.userRepo.FindAllUsers()
	if err != nil {
		return nil, err
	}

	res := make([]user.UserResponse, len(users))
	for i, u := range users {
		res[i] = user.UserResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			RoleName: u.Role.Name,
			Status:   u.Status,
		}
	}
	return res, nil
}

func (s *userMgmtService) UpdateUserStatus(userID uuid.UUID, status string) error {
	if status != "ACTIVE" && status != "SUSPENDED" {
		return errors.New("invalid status")
	}
	return s.userRepo.UpdateUserStatus(userID, status)
}

func generateRandomPassword(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buffer)[:length], nil
}
