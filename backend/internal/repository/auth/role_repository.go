package auth

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByID(id uint) (*domain.Role, error)
	FindPermissionsByRoleID(roleID uint) ([]domain.Permission, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByID(id uint) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindPermissionsByRoleID(roleID uint) ([]domain.Permission, error) {
	var rolePermissions []domain.RolePermission
	if err := r.db.Preload("Permission").Where("role_id = ?", roleID).Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	permissions := make([]domain.Permission, len(rolePermissions))
	for i, rp := range rolePermissions {
		permissions[i] = rp.Permission
	}
	return permissions, nil
}
