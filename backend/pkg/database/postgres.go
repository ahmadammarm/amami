package database

import (
	"fmt"
	"log"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
		cfg.DBTimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Create pg_trgm extension for GIN index search
	db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm;")

	// Auto Migration
	err = db.AutoMigrate(
		&domain.Role{},
		&domain.Permission{},
		&domain.RolePermission{},
		&domain.User{},
		&domain.Jamaah{},
		&domain.Fund{},
		&domain.Transaction{},
		&domain.ZakatDonation{},
		&domain.MustahikData{},
		&domain.ZakatDistribution{},
		&domain.QurbanPackage{},
		&domain.QurbanBooking{},
		&domain.QurbanAnimal{},
		&domain.Asset{},
		&domain.AssetLoan{},
		&domain.SystemSetting{},
		&domain.Memo{},
		&domain.Agenda{},
		&domain.AgendaDocumentation{},
		&domain.AuditLog{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}
