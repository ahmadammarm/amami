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
	if err := r.db.Select("id", "name").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (h *roleRepository) FindPermissionsByRoleID(roleID uint) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := h.db.Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Select("permissions.id", "permissions.code"). // Only select necessary fields
		Find(&permissions).Error
	
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
