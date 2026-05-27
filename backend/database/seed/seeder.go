package seed

import (
	"log"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	log.Println("Seeding database...")

	// 1. Seed Permissions
	permissions := []domain.Permission{
		{Code: "user:manage", Description: "Manage users and roles"},
		{Code: "ledger:read", Description: "Read financial reports"},
		{Code: "ledger:write", Description: "Create transactions"},
		{Code: "zakat:manage", Description: "Manage zakat collection and distribution"},
		{Code: "qurban:manage", Description: "Manage qurban packages and bookings"},
		{Code: "inventory:manage", Description: "Manage mosque assets"},
	}

	for _, p := range permissions {
		if err := db.FirstOrCreate(&p, domain.Permission{Code: p.Code}).Error; err != nil {
			return err
		}
	}

	// 2. Seed Roles
	roles := []domain.Role{
		{Name: "SUPER_ADMIN", Description: "Full system access"},
		{Name: "BENDAHARA", Description: "Financial management access"},
		{Name: "TAKMIR", Description: "General operational access"},
	}

	for i := range roles {
		if err := db.FirstOrCreate(&roles[i], domain.Role{Name: roles[i].Name}).Error; err != nil {
			return err
		}
	}

	// 3. Map Permissions to Roles (RBAC)
	var allPerms []domain.Permission
	db.Find(&allPerms)

	for _, p := range allPerms {
		rp := domain.RolePermission{RoleID: roles[0].ID, PermissionID: p.ID}
		db.FirstOrCreate(&rp, domain.RolePermission{RoleID: roles[0].ID, PermissionID: p.ID})
	}

	// 4. Seed Super Admin User
	hashedPassword, _ := utils.HashPassword("admin123")
	superAdmin := domain.User{
		Username:     "admin",
		Email:        "admin@amami.org",
		PasswordHash: hashedPassword,
		RoleID:       roles[0].ID,
		Status:       "ACTIVE",
	}

	if err := db.FirstOrCreate(&superAdmin, domain.User{Username: "admin"}).Error; err != nil {
		return err
	}

	log.Println("Seeding completed successfully!")
	return nil
}
