package user

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMgmtRepository interface {
	CreateUserWithProfile(user *domain.User, jamaah *domain.Jamaah) error
	FindAllUsers() ([]domain.User, error)
	UpdateUserStatus(userID uuid.UUID, status string) error
	FindByID(id uuid.UUID) (*domain.User, error)
}

type userMgmtRepository struct {
	db *gorm.DB
}

func NewUserMgmtRepository(db *gorm.DB) UserMgmtRepository {
	return &userMgmtRepository{db: db}
}

func (r *userMgmtRepository) CreateUserWithProfile(user *domain.User, jamaah *domain.Jamaah) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		
		jamaah.UserID = &user.ID
		if err := tx.Create(jamaah).Error; err != nil {
			return err
		}
		
		return nil
	})
}

func (r *userMgmtRepository) FindAllUsers() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Preload("Role").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userMgmtRepository) UpdateUserStatus(userID uuid.UUID, status string) error {
	return r.db.Model(&domain.User{}).Where("id = ?", userID).Update("status", status).Error
}

func (r *userMgmtRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload("Role").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
