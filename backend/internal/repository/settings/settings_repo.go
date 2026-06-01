package settings

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettingsRepository interface {
	UpsertSetting(key string, value string, isSecret bool) error
	GetSetting(key string) (*domain.SystemSetting, error)
	GetAllSettings() ([]domain.SystemSetting, error)
	Ping() error
	GetAuditLogs(page, limit int, userID uuid.UUID, roleID uint) ([]domain.AuditLog, int64, error)
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

func (r *settingsRepository) GetAuditLogs(page, limit int, userID uuid.UUID, roleID uint) ([]domain.AuditLog, int64, error) {
	var logs []domain.AuditLog
	var total int64

	// 1. Get Requester Role Name
	var requesterRole domain.Role
	if err := r.db.Select("name").First(&requesterRole, roleID).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Model(&domain.AuditLog{}).
		Joins("JOIN users ON users.id = audit_logs.user_id").
		Joins("JOIN roles ON roles.id = users.role_id")

	// Apply hierarchical filters based on requester role
	switch requesterRole.Name {
	case "SUPER_ADMIN":
		// No filter, see everything
	case "TAKMIR":
		// Sees own, BENDAHARA, and SEKRETARIS
		query = query.Where("roles.name IN ? OR audit_logs.user_id = ?", []string{"BENDAHARA", "SEKRETARIS"}, userID)
	case "BENDAHARA", "SEKRETARIS":
		// See each other and own
		query = query.Where("roles.name IN ? OR audit_logs.user_id = ?", []string{"BENDAHARA", "SEKRETARIS"}, userID)
	case "JAMAAH":
		// Only see JAMAAH's activities
		query = query.Where("roles.name = ?", "JAMAAH")
	default:
		// Default to own logs if role is unknown
		query = query.Where("audit_logs.user_id = ?", userID)
	}

	// Count total records for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username")
	}).
		Select("audit_logs.id", "audit_logs.user_id", "audit_logs.action", "audit_logs.entity", "audit_logs.created_at").
		Order("audit_logs.created_at DESC")
	
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
