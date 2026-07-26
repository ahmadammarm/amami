package user

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMgmtRepository interface {
	CreateUserWithProfile(user *domain.User, jamaah *domain.Jamaah) error
	FindAllUsers(page int, limit int) ([]domain.User, int64, error)
	UpdateUserStatus(userID uuid.UUID, status string) error
	DeleteUser(userID uuid.UUID) error
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

func (r *userMgmtRepository) FindAllUsers(page int, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	if err := r.db.Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Preload("Role").Preload("JamaahProfile").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userMgmtRepository) UpdateUserStatus(userID uuid.UUID, status string) error {
	return r.db.Model(&domain.User{}).Where("id = ?", userID).Update("status", status).Error
}

func (r *userMgmtRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload("Role").Preload("JamaahProfile").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userMgmtRepository) DeleteUser(userID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// First delete Jamaah Profile if it exists
		if err := tx.Where("user_id = ?", userID).Delete(&domain.Jamaah{}).Error; err != nil {
			return err
		}
		// Then delete the User
		if err := tx.Where("id = ?", userID).Delete(&domain.User{}).Error; err != nil {
			return err
		}
		return nil
	})
}
