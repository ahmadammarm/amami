package settings

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettingsRepository interface {
	// UpsertSettings writes every setting in one statement, so a save is
	// all-or-nothing and concurrent saves never interleave key by key.
	UpsertSettings(settings []domain.SystemSetting) error
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

// UpsertSettings issues a single INSERT ... ON CONFLICT (key) DO UPDATE.
// Callers pass the keys in a fixed order, so two concurrent saves lock the
// rows in the same order and cannot deadlock.
func (r *settingsRepository) UpsertSettings(settings []domain.SystemSetting) error {
	if len(settings) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "is_secret"}),
	}).Create(&settings).Error
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
