package auth

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id uuid.UUID) (*domain.User, error)
	FindByUsernameOrEmail(identifier string) (*domain.User, error)
	Create(user *domain.User) error
	UpdateCredentials(id uuid.UUID, passwordHash string, status domain.UserStatus) error
	CreateAuditLog(log domain.AuditLog) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.Select("id", "username", "email", "role_id", "status").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("JamaahProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id", "full_name")
		}).
		First(&user, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsernameOrEmail(identifier string) (*domain.User, error) {
	var user domain.User
	err := r.db.Select("id", "username", "email", "password_hash", "role_id", "status").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Where("username = ? OR email = ?", identifier, identifier).
		First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// UpdateCredentials updates only the password and status columns, so partially
// loaded users (see FindByID) never overwrite other fields with zero values.
func (r *userRepository) UpdateCredentials(id uuid.UUID, passwordHash string, status domain.UserStatus) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password_hash": passwordHash,
		"status":        status,
	}).Error
}

func (r *userRepository) CreateAuditLog(log domain.AuditLog) error {
	return r.db.Create(&log).Error
}
