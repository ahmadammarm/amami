package user

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMgmtRepository interface {
	CreateUserWithProfile(user *domain.User, jamaah *domain.Jamaah) error
	FindAllUsers(page int, limit int) ([]domain.User, int64, error)
	UpdateUserStatus(userID uuid.UUID, status domain.UserStatus) error
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
		return tx.Create(jamaah).Error
	})
}

func (r *userMgmtRepository) FindAllUsers(page int, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	if err := r.db.Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Preload("Role").Preload("JamaahProfile").
		Order("created_at desc").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userMgmtRepository) UpdateUserStatus(userID uuid.UUID, status domain.UserStatus) error {
	res := r.db.Model(&domain.User{}).Where("id = ?", userID).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userMgmtRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload("Role").Preload("JamaahProfile").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser soft-deletes the account and its linked jamaah profile.
func (r *userMgmtRepository) DeleteUser(userID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&domain.Jamaah{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", userID).Delete(&domain.User{}).Error
	})
}
