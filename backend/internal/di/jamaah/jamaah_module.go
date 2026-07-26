package jamaah

import (
	"github.com/ahmadammarm/amami/backend/internal/handler/jamaah"
	jamaahRepo "github.com/ahmadammarm/amami/backend/internal/repository/jamaah"
	jamaahSvc "github.com/ahmadammarm/amami/backend/internal/service/jamaah"
	"gorm.io/gorm"
)

type Module struct {
	Handler *jamaah.JamaahHandler
}

func NewModule(db *gorm.DB) *Module {
	repo := jamaahRepo.NewJamaahRepository(db)
	svc := jamaahSvc.NewJamaahService(repo)
	handler := jamaah.NewJamaahHandler(svc)

	return &Module{
		Handler: handler,
	}
}
