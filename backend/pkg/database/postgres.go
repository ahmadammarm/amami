package database

import (
	"fmt"
	"log"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// partialUniqueIndexes were plain unique indexes before soft deletes were added.
// AutoMigrate never alters an existing index, so the old full-table versions are
// dropped first and recreated by AutoMigrate with `WHERE deleted_at IS NULL`.
var partialUniqueIndexes = []string{
	"idx_users_username",
	"idx_users_email",
	"idx_assets_sku",
	"idx_qurban_animals_tag_number",
}

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

	// TranslateError turns driver errors into gorm.ErrDuplicatedKey and friends.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Create pg_trgm extension for GIN index search
	db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm;")

	if err := dropNonPartialUniqueIndexes(db); err != nil {
		log.Fatalf("failed to prepare indexes: %v", err)
	}

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

	// Transaction descriptions used to live in metadata JSON; copy them into
	// the dedicated column. Idempotent: only rows with an empty column match.
	if err := db.Exec(`UPDATE transactions SET description = metadata->>'description'
		WHERE COALESCE(description, '') = '' AND metadata->>'description' IS NOT NULL`).Error; err != nil {
		log.Fatalf("failed to backfill transaction descriptions: %v", err)
	}

	return db
}

func dropNonPartialUniqueIndexes(db *gorm.DB) error {
	for _, name := range partialUniqueIndexes {
		var count int64
		err := db.Raw(
			"SELECT COUNT(*) FROM pg_indexes WHERE indexname = ? AND indexdef NOT ILIKE '%WHERE%'",
			name,
		).Scan(&count).Error
		if err != nil {
			return err
		}
		if count == 0 {
			continue
		}
		if err := db.Exec(fmt.Sprintf("DROP INDEX IF EXISTS %q", name)).Error; err != nil {
			return err
		}
	}
	return nil
}
