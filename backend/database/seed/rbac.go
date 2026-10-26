package seed

import (
	"fmt"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var roleDescriptions = map[string]string{
	domain.RoleSuperAdmin: "Full system governance and management access",
	domain.RoleBendahara:  "Financial and Zakat management access",
	domain.RoleTakmir:     "Broad operational management access (All except User/System management)",
	domain.RoleSekretaris: "Administrative and documentation access",
	domain.RoleJamaah:     "Community member access (Personal data only)",
}

// SyncRBAC makes the permissions, roles and role_permissions tables match
// domain.PermissionDescriptions and domain.RolePermissions exactly. It is
// idempotent and runs on every server start, so grant changes in code take
// effect after a restart without re-running the sample-data seeder.
func SyncRBAC(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		permIDs, err := syncPermissions(tx)
		if err != nil {
			return err
		}

		for roleName, codes := range domain.RolePermissions {
			role := domain.Role{Name: roleName, Description: roleDescriptions[roleName]}
			if err := tx.Where(domain.Role{Name: roleName}).
				Assign(domain.Role{Description: role.Description}).
				FirstOrCreate(&role).Error; err != nil {
				return fmt.Errorf("sync role %s: %w", roleName, err)
			}
			if err := syncRoleGrants(tx, role.ID, codes, permIDs); err != nil {
				return fmt.Errorf("sync grants for %s: %w", roleName, err)
			}
		}
		return nil
	})
}

func syncPermissions(tx *gorm.DB) (map[domain.PermissionCode]uint, error) {
	ids := make(map[domain.PermissionCode]uint, len(domain.PermissionDescriptions))
	for code, description := range domain.PermissionDescriptions {
		perm := domain.Permission{Code: string(code)}
		if err := tx.Where(domain.Permission{Code: string(code)}).
			Assign(domain.Permission{Description: description}).
			FirstOrCreate(&perm).Error; err != nil {
			return nil, fmt.Errorf("sync permission %s: %w", code, err)
		}
		ids[code] = perm.ID
	}
	return ids, nil
}

func syncRoleGrants(tx *gorm.DB, roleID uint, codes []domain.PermissionCode, permIDs map[domain.PermissionCode]uint) error {
	wanted := make([]uint, 0, len(codes))
	for _, code := range codes {
		id, ok := permIDs[code]
		if !ok {
			return fmt.Errorf("unknown permission %s", code)
		}
		wanted = append(wanted, id)
	}

	if err := tx.Where("role_id = ? AND permission_id NOT IN ?", roleID, wanted).
		Delete(&domain.RolePermission{}).Error; err != nil {
		return err
	}

	grants := make([]domain.RolePermission, 0, len(wanted))
	for _, permID := range wanted {
		grants = append(grants, domain.RolePermission{RoleID: roleID, PermissionID: permID})
	}
	return tx.Clauses(clause.OnConflict{DoNothing: true}).Omit(clause.Associations).Create(&grants).Error
}
