package di

import (
	authHandler "github.com/ahmadammarm/amami/backend/internal/handler/auth"
	authRepo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	authSvc "github.com/ahmadammarm/amami/backend/internal/service/auth"
	dashboardModule "github.com/ahmadammarm/amami/backend/internal/di/dashboard"
	financeModule "github.com/ahmadammarm/amami/backend/internal/di/finance"
	jamaahModule "github.com/ahmadammarm/amami/backend/internal/di/jamaah"
	settingsModule "github.com/ahmadammarm/amami/backend/internal/di/settings"
	userModule "github.com/ahmadammarm/amami/backend/internal/di/user"
	zakatModule "github.com/ahmadammarm/amami/backend/internal/di/zakat"
	agendaModule "github.com/ahmadammarm/amami/backend/internal/di/agenda"
	inventoryModule "github.com/ahmadammarm/amami/backend/internal/di/inventory"
	qurbanModule "github.com/ahmadammarm/amami/backend/internal/di/qurban"
	"github.com/ahmadammarm/amami/backend/pkg/config"
	"github.com/ahmadammarm/amami/backend/pkg/database"
	"gorm.io/gorm"
)

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

	// Base/Common Repositories
	roleRepo := authRepo.NewRoleRepository(db)

	// Auth Module
	userRepo := authRepo.NewUserRepository(db)
	authService := authSvc.NewAuthService(userRepo)
	authH := authHandler.NewAuthHandler(authService)

	// User Management Module
	userMod := userModule.NewModule(db, roleRepo)

	// Settings Module
	settingsMod := settingsModule.NewModule(db)

	// Finance Module
	financeMod := financeModule.NewModule(db)

	// Dashboard Module
	dashboardMod := dashboardModule.NewModule(db)

	// Jamaah Module
	jamaahMod := jamaahModule.NewModule(db)

	// Zakat Module
	zakatMod := zakatModule.NewModule(db)

	// Agenda Module
	agendaMod := agendaModule.NewModule(db)

	// Inventory Module
	inventoryMod := inventoryModule.NewModule(db)

	// Qurban Module
	qurbanMod := qurbanModule.NewModule(db)

	return &Config{
		AuthHandler:      authH,
		UserHandler:      *userMod,
		SettingsHandler:  *settingsMod,
		FinanceHandler:   *financeMod,
		DashboardHandler: *dashboardMod,
		JamaahHandler:    *jamaahMod,
		ZakatHandler:     *zakatMod,
		AgendaHandler:    *agendaMod,
		InventoryHandler: *inventoryMod,
		QurbanHandler:    *qurbanMod,
		RoleRepo:         roleRepo,
		DB:               db,
		AppConfig:        cfg,
	}, nil
}
