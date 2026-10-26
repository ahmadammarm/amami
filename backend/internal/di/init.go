package di

import (
	"fmt"

	"github.com/ahmadammarm/amami/backend/database/seed"
	agendaModule "github.com/ahmadammarm/amami/backend/internal/di/agenda"
	dashboardModule "github.com/ahmadammarm/amami/backend/internal/di/dashboard"
	financeModule "github.com/ahmadammarm/amami/backend/internal/di/finance"
	inventoryModule "github.com/ahmadammarm/amami/backend/internal/di/inventory"
	jamaahModule "github.com/ahmadammarm/amami/backend/internal/di/jamaah"
	qurbanModule "github.com/ahmadammarm/amami/backend/internal/di/qurban"
	settingsModule "github.com/ahmadammarm/amami/backend/internal/di/settings"
	userModule "github.com/ahmadammarm/amami/backend/internal/di/user"
	zakatModule "github.com/ahmadammarm/amami/backend/internal/di/zakat"
	authHandler "github.com/ahmadammarm/amami/backend/internal/handler/auth"
	authRepo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	authSvc "github.com/ahmadammarm/amami/backend/internal/service/auth"
	"github.com/ahmadammarm/amami/backend/pkg/config"
	"github.com/ahmadammarm/amami/backend/pkg/database"
	"github.com/ahmadammarm/amami/backend/pkg/storage"
	"gorm.io/gorm"
)

// UploadURLPrefix is the route under which stored uploads are served.
const UploadURLPrefix = "/uploads"

type Config struct {
	AuthHandler      authHandler.AuthHandler
	UserHandler      userModule.Module
	SettingsHandler  settingsModule.Module
	FinanceHandler   financeModule.Module
	DashboardHandler dashboardModule.Module
	JamaahHandler    jamaahModule.Module
	ZakatHandler     zakatModule.Module
	AgendaHandler    agendaModule.Module
	InventoryHandler inventoryModule.Module
	QurbanHandler    qurbanModule.Module
	RoleRepo         authRepo.RoleRepository
	DB               *gorm.DB
	AppConfig        *config.Config
}

func Initialize() (*Config, error) {
	cfg := config.LoadConfig()
	db := database.NewDatabase(cfg)

	// Keep roles and grants in line with domain.RolePermissions on every start.
	if err := seed.SyncRBAC(db); err != nil {
		return nil, fmt.Errorf("sync RBAC: %w", err)
	}

	// Base/Common Repositories
	roleRepo := authRepo.NewRoleRepository(db)
	uploads := storage.NewLocalStorage(cfg.UploadDir, UploadURLPrefix)

	// Auth Module
	userRepo := authRepo.NewUserRepository(db)
	authService := authSvc.NewAuthService(userRepo, roleRepo)
	authH := authHandler.NewAuthHandler(authService)

	return &Config{
		AuthHandler:      authH,
		UserHandler:      *userModule.NewModule(db, roleRepo),
		SettingsHandler:  *settingsModule.NewModule(db),
		FinanceHandler:   *financeModule.NewModule(db),
		DashboardHandler: *dashboardModule.NewModule(db),
		JamaahHandler:    *jamaahModule.NewModule(db),
		ZakatHandler:     *zakatModule.NewModule(db),
		AgendaHandler:    *agendaModule.NewModule(db, uploads, cfg.MaxUploadBytes),
		InventoryHandler: *inventoryModule.NewModule(db),
		QurbanHandler:    *qurbanModule.NewModule(db),
		RoleRepo:         roleRepo,
		DB:               db,
		AppConfig:        cfg,
	}, nil
}
