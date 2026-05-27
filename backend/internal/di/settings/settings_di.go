package settings

import (
	settingsHandler "github.com/ahmadammarm/amami/backend/internal/handler/settings"
	settingsRepo "github.com/ahmadammarm/amami/backend/internal/repository/settings"
	settingsSvc "github.com/ahmadammarm/amami/backend/internal/service/settings"
	"gorm.io/gorm"
)

type Module struct {
	Handler settingsHandler.SettingsHandler
}

func NewModule(db *gorm.DB) *Module {
	repo := settingsRepo.NewSettingsRepository(db)
	svc := settingsSvc.NewSettingsService(repo)
	h := settingsHandler.NewSettingsHandler(svc)
	
	return &Module{Handler: h}
}
