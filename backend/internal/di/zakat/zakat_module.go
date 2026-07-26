package zakat

import (
	"github.com/ahmadammarm/amami/backend/internal/handler/zakat"
	zakatRepo "github.com/ahmadammarm/amami/backend/internal/repository/zakat"
	zakatSvc "github.com/ahmadammarm/amami/backend/internal/service/zakat"
	"gorm.io/gorm"
)

type Module struct {
	Handler *zakat.ZakatHandler
}

func NewModule(db *gorm.DB) *Module {
	repo := zakatRepo.NewZakatRepository(db)
	svc := zakatSvc.NewZakatService(repo)
	handler := zakat.NewZakatHandler(svc)

	return &Module{
		Handler: handler,
	}
}
