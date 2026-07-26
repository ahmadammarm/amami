package finance

import (
	"github.com/ahmadammarm/amami/backend/internal/handler/finance"
	financeRepo "github.com/ahmadammarm/amami/backend/internal/repository/finance"
	financeSvc "github.com/ahmadammarm/amami/backend/internal/service/finance"
	"gorm.io/gorm"
)

type Module struct {
	Handler *finance.FinanceHandler
}

func NewModule(db *gorm.DB) *Module {
	repo := financeRepo.NewFinanceRepository(db)
	svc := financeSvc.NewFinanceService(repo)
	handler := finance.NewFinanceHandler(svc)

	return &Module{
		Handler: handler,
	}
}
