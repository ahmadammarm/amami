package auth

import (
	"sync"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByID(id uint) (*domain.Role, error)
	FindPermissionsByRoleID(roleID uint) ([]domain.Permission, error)
}

type roleRepository struct {
	db              *gorm.DB
	permissionCache sync.Map
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) FindByID(id uint) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.Select("id", "name").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindPermissionsByRoleID(roleID uint) ([]domain.Permission, error) {
	// Check cache first
	if val, ok := r.permissionCache.Load(roleID); ok {
		return val.([]domain.Permission), nil
	}

	// If not in cache, query the database
	var permissions []domain.Permission
	err := r.db.Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Select("permissions.id", "permissions.code").
		Find(&permissions).Error
	
	if err != nil {
		return nil, err
	}

	// Store in cache
	r.permissionCache.Store(roleID, permissions)

	return permissions, nil
}
