package user

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/dto/user"
	authRepo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	userRepo "github.com/ahmadammarm/amami/backend/internal/repository/user"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserMgmtService interface {
	InviteUser(req user.InviteUserRequest) (*user.InviteUserResponse, error)
	GetAllUsers(page, limit int) (utils.Page[user.UserResponse], error)
	GetRoles() ([]user.RoleResponse, error)
	UpdateUserStatus(userID uuid.UUID, status domain.UserStatus) error
	DeleteUser(actorID, userID uuid.UUID) error
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

func (s *userMgmtService) InviteUser(req user.InviteUserRequest) (*user.InviteUserResponse, error) {
	role, err := s.roleRepo.FindByID(req.RoleID)
	if err != nil {
		return nil, apperror.Validation("invalid role id")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.Internal("failed to hash password", err)
	}

	userModel := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		RoleID:       req.RoleID,
		Role:         *role,
		Status:       domain.UserPendingPasswordChange,
	}
	jamaah := &domain.Jamaah{FullName: req.FullName}

	if err := s.userRepo.CreateUserWithProfile(userModel, jamaah); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("username or email already exists")
		}
		return nil, apperror.Internal("failed to create user", err)
	}
	userModel.JamaahProfile = jamaah

	zap.L().Info("User invited successfully", zap.String("username", userModel.Username))

	// The temporary password is returned until invitation emails are sent over SMTP.
	return &user.InviteUserResponse{
		User:              user.ToUserResponse(userModel),
		TemporaryPassword: req.Password,
	}, nil
}

func (s *userMgmtService) GetAllUsers(page, limit int) (utils.Page[user.UserResponse], error) {
	users, total, err := s.userRepo.FindAllUsers(page, limit)
	if err != nil {
		return utils.Page[user.UserResponse]{}, apperror.Internal("failed to fetch users", err)
	}

	res := make([]user.UserResponse, 0, len(users))
	for i := range users {
		res = append(res, user.ToUserResponse(&users[i]))
	}
	return utils.NewPage(res, total, page, limit), nil
}

func (s *userMgmtService) GetRoles() ([]user.RoleResponse, error) {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, apperror.Internal("failed to fetch roles", err)
	}

	res := make([]user.RoleResponse, 0, len(roles))
	for _, r := range roles {
		res = append(res, user.RoleResponse{ID: r.ID, Name: r.Name, Description: r.Description})
	}
	return res, nil
}

func (s *userMgmtService) UpdateUserStatus(userID uuid.UUID, status domain.UserStatus) error {
	if err := s.userRepo.UpdateUserStatus(userID, status); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("user not found")
		}
		return apperror.Internal("failed to update user status", err)
	}
	return nil
}

func (s *userMgmtService) DeleteUser(actorID, userID uuid.UUID) error {
	if actorID == userID {
		return apperror.Forbidden("you cannot delete your own account")
	}

	target, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("user not found")
		}
		return apperror.Internal("failed to load user", err)
	}

	if target.Role.Name == domain.RoleSuperAdmin {
		return apperror.Forbidden("cannot delete a SUPER_ADMIN")
	}

	if err := s.userRepo.DeleteUser(userID); err != nil {
		return apperror.Internal("failed to delete user", err)
	}
	return nil
}
