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
	Ping() error
	GetAuditLogs(page, limit int) ([]domain.AuditLog, int64, error)
	CreateAuditLog(log domain.AuditLog) error
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
	if err := r.db.Select("id", "key", "value", "is_secret").Where("key = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *settingsRepository) GetAllSettings() ([]domain.SystemSetting, error) {
	var settings []domain.SystemSetting
	if err := r.db.Select("id", "key", "value", "is_secret").Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *settingsRepository) Ping() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (r *settingsRepository) GetAuditLogs(page, limit int) ([]domain.AuditLog, int64, error) {
	var logs []domain.AuditLog
	var total int64

	// Count total records for pagination
	if err := r.db.Model(&domain.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username")
	}).
		Select("id", "user_id", "action", "entity", "created_at").
		Order("created_at DESC")
	
	if limit > 0 {
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (r *settingsRepository) CreateAuditLog(log domain.AuditLog) error {
	return r.db.Create(&log).Error
}
