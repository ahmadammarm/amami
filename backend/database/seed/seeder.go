package seed

import (
	"log"
	"time"

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
		{Code: "jamaah:manage", Description: "Manage community database and scoring"},
		{Code: "agenda:manage", Description: "Manage mosque events and documentation"},
		{Code: "settings:manage", Description: "Manage mosque profile and system settings"},
	}

	for _, p := range permissions {
		if err := db.FirstOrCreate(&p, domain.Permission{Code: p.Code}).Error; err != nil {
			return err
		}
	}

	// 2. Seed Roles
	roles := []domain.Role{
		{Name: "SUPER_ADMIN", Description: "Full system governance and management access"},
		{Name: "BENDAHARA", Description: "Financial and Zakat management access"},
		{Name: "TAKMIR", Description: "Broad operational management access (All except User/System management)"},
		{Name: "SEKRETARIS", Description: "Administrative and documentation access"},
		{Name: "JAMAAH", Description: "Community member access (Personal data only)"},
	}

	for i := range roles {
		if err := db.FirstOrCreate(&roles[i], domain.Role{Name: roles[i].Name}).Error; err != nil {
			return err
		}
	}

	// 3. Map Permissions to Roles (RBAC)
	
	// Helper to find permission ID by code
	getPermID := func(code string) uint {
		var p domain.Permission
		db.Where("code = ?", code).First(&p)
		return p.ID
	}

	// SUPER_ADMIN: Gets everything
	var allPerms []domain.Permission
	db.Find(&allPerms)
	for _, p := range allPerms {
		db.FirstOrCreate(&domain.RolePermission{RoleID: roles[0].ID, PermissionID: p.ID})
	}

	// BENDAHARA: Finance, Zakat, and reading Jamaah
	bendaharaPerms := []string{"ledger:read", "ledger:write", "zakat:manage", "jamaah:manage"}
	for _, code := range bendaharaPerms {
		db.FirstOrCreate(&domain.RolePermission{RoleID: roles[1].ID, PermissionID: getPermID(code)})
	}

	// TAKMIR: Broad access (Everything except user:manage and settings:manage)
	takmirPerms := []string{"ledger:read", "zakat:manage", "qurban:manage", "inventory:manage", "jamaah:manage", "agenda:manage"}
	for _, code := range takmirPerms {
		db.FirstOrCreate(&domain.RolePermission{RoleID: roles[2].ID, PermissionID: getPermID(code)})
	}

	// SEKRETARIS: Administration, Agenda, and reading reports
	sekretarisPerms := []string{"ledger:read", "jamaah:manage", "agenda:manage", "settings:manage"}
	for _, code := range sekretarisPerms {
		db.FirstOrCreate(&domain.RolePermission{RoleID: roles[3].ID, PermissionID: getPermID(code)})
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

	// 5. Seed System Settings
	defaultSettings := []domain.SystemSetting{
		{Key: "mosque_name", Value: "Masjid Agung Al-Hikmah"},
		{Key: "mosque_address", Value: "Jl. Merdeka No. 1, Kota Amanah"},
		{Key: "mosque_phone", Value: "021-1234567"},
		{Key: "mosque_logo", Value: "/icons.svg"},
		{Key: "legal_yayasan_id", Value: "YAYASAN-AH-001"},
		{Key: "smtp_host", Value: "smtp.example.com"},
		{Key: "smtp_port", Value: "587"},
		{Key: "smtp_user", Value: "admin@amami.org"},
		{Key: "smtp_from", Value: "no-reply@amami.org"},
	}

	for _, s := range defaultSettings {
		if err := db.FirstOrCreate(&s, domain.SystemSetting{Key: s.Key}).Error; err != nil {
			return err
		}
	}

	// 5.5 Seed Default Funds
	defaultFunds := []domain.Fund{
		{Name: "Kas Operasional Masjid", Code: "KAS-OPR", CurrentBalance: 0},
		{Name: "Kas Pembangunan", Code: "KAS-PEM", CurrentBalance: 0},
		{Name: "Kas Anak Yatim", Code: "KAS-YAT", CurrentBalance: 0},
	}
	for _, f := range defaultFunds {
		if err := db.FirstOrCreate(&f, domain.Fund{Code: f.Code}).Error; err != nil {
			return err
		}
	}

	// 6. Seed Audit Logs
	for i := 1; i <= 25; i++ {
		action := "UPDATE"
		if i%3 == 0 { action = "LOGIN" }
		if i%5 == 0 { action = "CREATE" }
		
		log := domain.AuditLog{
			UserID:    superAdmin.ID,
			Action:    action,
			Entity:    "System Configuration",
			CreatedAt: time.Now().Add(time.Duration(-i) * time.Hour),
		}
		db.Create(&log)
	}

	log.Println("Seeding completed successfully!")
	return nil
}
