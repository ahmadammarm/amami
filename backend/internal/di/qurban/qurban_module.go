package qurban

import (
	qurbanHandler "github.com/ahmadammarm/amami/backend/internal/handler/qurban"
	qurbanRepo "github.com/ahmadammarm/amami/backend/internal/repository/qurban"
	qurbanSvc "github.com/ahmadammarm/amami/backend/internal/service/qurban"

	"gorm.io/gorm"
)

type Module struct {
	Handler *qurbanHandler.QurbanHandler
}

func NewModule(db *gorm.DB) *Module {
	repo := qurbanRepo.NewQurbanRepository(db)
	svc := qurbanSvc.NewQurbanService(repo)
	handler := qurbanHandler.NewQurbanHandler(svc)

	return &Module{
		Handler: handler,
	}
}
