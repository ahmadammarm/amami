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

// models lists every table AutoMigrate manages; a new model must be added here.
var models = []any{
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
}

// NewDatabase connects and migrates, exiting the process on failure.
// It is the entry point for the server; tests use Open and Migrate directly.
func NewDatabase(cfg *config.Config) *gorm.DB {
	db, err := Open(DSN(cfg))
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if err := Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	return db
}

// DSN builds the PostgreSQL connection string from the configuration.
func DSN(cfg *config.Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
		cfg.DBTimeZone,
	)
}

// Open connects to PostgreSQL. TranslateError turns driver errors into
// gorm.ErrDuplicatedKey and gorm.ErrForeignKeyViolated.
func Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
}

// Migrate brings the schema up to date. Every step is idempotent.
func Migrate(db *gorm.DB) error {
	// Trigram operators for the fuzzy jamaah name index.
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error; err != nil {
		return fmt.Errorf("create pg_trgm extension: %w", err)
	}

	if err := dropNonPartialUniqueIndexes(db); err != nil {
		return fmt.Errorf("prepare unique indexes: %w", err)
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto-migrate: %w", err)
	}

	// Transaction descriptions used to live in metadata JSON; copy them into
	// the dedicated column. Only rows with an empty column match.
	if err := db.Exec(`UPDATE transactions SET description = metadata->>'description'
		WHERE COALESCE(description, '') = '' AND metadata->>'description' IS NOT NULL`).Error; err != nil {
		return fmt.Errorf("backfill transaction descriptions: %w", err)
	}
	return nil
}

func dropNonPartialUniqueIndexes(db *gorm.DB) error {
	for _, name := range partialUniqueIndexes {
		var count int64
		err := db.Raw(
			`SELECT COUNT(*) FROM pg_indexes
			WHERE schemaname = current_schema() AND indexname = ? AND indexdef NOT ILIKE '%WHERE%'`,
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
