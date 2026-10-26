package auth

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/dto/auth"
	"github.com/ahmadammarm/amami/backend/internal/dto/user"
	repo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthService interface {
	Login(req auth.LoginRequest) (*auth.LoginResponse, error)
	GetMe(userID uuid.UUID) (*user.MeResponse, error)
	ChangePassword(userID uuid.UUID, req auth.ChangePasswordRequest) (*auth.LoginResponse, error)
}

type authService struct {
	userRepo repo.UserRepository
	roleRepo repo.RoleRepository
}

func NewAuthService(userRepo repo.UserRepository, roleRepo repo.RoleRepository) AuthService {
	return &authService{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *authService) Login(req auth.LoginRequest) (*auth.LoginResponse, error) {
	u, err := s.userRepo.FindByUsernameOrEmail(req.Email)
	if err != nil {
		return nil, apperror.Unauthorized("invalid credentials")
	}

	if u.Status != domain.UserActive && u.Status != domain.UserPendingPasswordChange {
		return nil, apperror.Unauthorized("account is suspended or inactive")
	}

	if !utils.CheckPasswordHash(req.Password, u.PasswordHash) {
		return nil, apperror.Unauthorized("invalid credentials")
	}

	scope := utils.ScopeFullAccess
	requiresChange := u.Status == domain.UserPendingPasswordChange
	if requiresChange {
		scope = utils.ScopePasswordReset
	}

	token, err := utils.GenerateToken(u.ID, u.RoleID, scope)
	if err != nil {
		return nil, apperror.Internal("failed to generate token", err)
	}

	s.recordAudit(u.ID, "LOGIN", "User Session")

	return &auth.LoginResponse{
		Token:                  token,
		RequiresPasswordChange: requiresChange,
	}, nil
}

func (s *authService) ChangePassword(userID uuid.UUID, req auth.ChangePasswordRequest) (*auth.LoginResponse, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperror.NotFound("user not found")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, apperror.Internal("failed to hash password", err)
	}

	if err := s.userRepo.UpdateCredentials(u.ID, hashedPassword, domain.UserActive); err != nil {
		return nil, apperror.Internal("failed to update user", err)
	}

	token, err := utils.GenerateToken(u.ID, u.RoleID, utils.ScopeFullAccess)
	if err != nil {
		return nil, apperror.Internal("failed to generate token", err)
	}

	s.recordAudit(u.ID, "PASSWORD_CHANGE", "User Account")

	return &auth.LoginResponse{
		Token:                  token,
		RequiresPasswordChange: false,
	}, nil
}

func (s *authService) GetMe(userID uuid.UUID) (*user.MeResponse, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperror.NotFound("user not found")
	}

	permissions, err := s.roleRepo.FindPermissionsByRoleID(u.RoleID)
	if err != nil {
		return nil, apperror.Internal("failed to load permissions", err)
	}

	codes := make([]domain.PermissionCode, 0, len(permissions))
	for _, p := range permissions {
		codes = append(codes, domain.PermissionCode(p.Code))
	}

	return &user.MeResponse{
		UserResponse: user.ToUserResponse(u),
		Permissions:  codes,
	}, nil
}

// recordAudit is best effort: a failed audit write must not fail the login itself.
func (s *authService) recordAudit(userID uuid.UUID, action, entity string) {
	if err := s.userRepo.CreateAuditLog(domain.AuditLog{UserID: userID, Action: action, Entity: entity}); err != nil {
		zap.L().Warn("failed to write audit log", zap.String("action", action), zap.Error(err))
	}
}
