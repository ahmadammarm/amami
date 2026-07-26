package dashboard

import (
	"github.com/ahmadammarm/amami/backend/internal/handler/dashboard"
	dashRepo "github.com/ahmadammarm/amami/backend/internal/repository/dashboard"
	dashSvc "github.com/ahmadammarm/amami/backend/internal/service/dashboard"
	"gorm.io/gorm"
)

type Module struct {
	Handler *dashboard.DashboardHandler
}

func NewModule(db *gorm.DB) *Module {
	repo := dashRepo.NewDashboardRepository(db)
	svc := dashSvc.NewDashboardService(repo)
	handler := dashboard.NewDashboardHandler(svc)

	return &Module{
		Handler: handler,
	}
}
