package settings

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettingsRepository interface {
	UpsertSetting(key string, value string, isSecret bool) error
	GetSetting(key string) (*domain.SystemSetting, error)
	GetAllSettings() ([]domain.SystemSetting, error)
}

type settingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) SettingsRepository {
	return &settingsRepository{db: db}
}

func (r *settingsRepository) UpsertSetting(key string, value string, isSecret bool) error {
	setting := domain.SystemSetting{
		Key:      key,
		Value:    value,
		IsSecret: isSecret,
	}
	
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "is_secret"}),
	}).Create(&setting).Error
}

func (r *settingsRepository) GetSetting(key string) (*domain.SystemSetting, error) {
	var setting domain.SystemSetting
	if err := r.db.Where("key = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *settingsRepository) GetAllSettings() ([]domain.SystemSetting, error) {
	var settings []domain.SystemSetting
	if err := r.db.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}
