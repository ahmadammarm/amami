package logistics

import (
	logisticsHandler "github.com/ahmadammarm/amami/backend/internal/handler/logistics"
	logisticsRepo "github.com/ahmadammarm/amami/backend/internal/repository/logistics"
	logisticsSvc "github.com/ahmadammarm/amami/backend/internal/service/logistics"
	"gorm.io/gorm"
)

type Module struct {
	Handler *logisticsHandler.LogisticsHandler
}

func NewModule(db *gorm.DB) *Module {
	repo := logisticsRepo.NewLogisticsRepository(db)
	svc := logisticsSvc.NewLogisticsService(repo)
	handler := logisticsHandler.NewLogisticsHandler(svc)

	return &Module{
		Handler: handler,
	}
}
