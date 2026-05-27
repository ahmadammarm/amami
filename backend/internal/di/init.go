package di

import (
	authHandler "github.com/ahmadammarm/amami/backend/internal/handler/auth"
	authRepo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	authSvc "github.com/ahmadammarm/amami/backend/internal/service/auth"
	userModule "github.com/ahmadammarm/amami/backend/internal/di/user"
	settingsModule "github.com/ahmadammarm/amami/backend/internal/di/settings"
	"github.com/ahmadammarm/amami/backend/pkg/config"
	"github.com/ahmadammarm/amami/backend/pkg/database"
	"gorm.io/gorm"
)

type Config struct {
	AuthHandler     authHandler.AuthHandler
	UserHandler     userModule.Module
	SettingsHandler settingsModule.Module
	RoleRepo        authRepo.RoleRepository
	DB              *gorm.DB
	AppConfig       *config.Config
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

	return &Config{
		AuthHandler:     authH,
		UserHandler:     *userMod,
		SettingsHandler: *settingsMod,
		RoleRepo:        roleRepo,
		DB:              db,
		AppConfig:       cfg,
	}, nil
}
