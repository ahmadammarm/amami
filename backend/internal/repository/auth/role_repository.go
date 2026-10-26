package auth

import (
	"sync"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindAll() ([]domain.Role, error)
	FindByID(id uint) (*domain.Role, error)
	FindPermissionsByRoleID(roleID uint) ([]domain.Permission, error)
}

type roleRepository struct {
	db *gorm.DB
	// permissionCache lives for the process lifetime; grants are synced on
	// startup (seed.SyncRBAC), so a restart picks up any change.
	permissionCache sync.Map
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) FindAll() ([]domain.Role, error) {
	var roles []domain.Role
	if err := r.db.Order("id").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) FindByID(id uint) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.Select("id", "name").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindPermissionsByRoleID(roleID uint) ([]domain.Permission, error) {
	if val, ok := r.permissionCache.Load(roleID); ok {
		return val.([]domain.Permission), nil
	}

	var permissions []domain.Permission
	err := r.db.Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Select("permissions.id", "permissions.code").
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	r.permissionCache.Store(roleID, permissions)
	return permissions, nil
}
