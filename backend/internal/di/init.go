package di

import (
	authHandler "github.com/ahmadammarm/amami/backend/internal/handler/auth"
	authRepo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	authSvc "github.com/ahmadammarm/amami/backend/internal/service/auth"
	"github.com/ahmadammarm/amami/backend/pkg/config"
	"github.com/ahmadammarm/amami/backend/pkg/database"
	"gorm.io/gorm"
)

type Config struct {
	AuthHandler authHandler.AuthHandler
	RoleRepo    authRepo.RoleRepository
	DB          *gorm.DB
	AppConfig   *config.Config
}

func Initialize() (*Config, error) {
	cfg := config.LoadConfig()
	db := database.NewDatabase(cfg)

	// Auth Module
	userRepo := authRepo.NewUserRepository(db)
	roleRepo := authRepo.NewRoleRepository(db)
	
	authService := authSvc.NewAuthService(userRepo)
	authH := authHandler.NewAuthHandler(authService)

	return &Config{
		AuthHandler: authH,
		RoleRepo:    roleRepo,
		DB:          db,
		AppConfig:   cfg,
	}, nil
}
