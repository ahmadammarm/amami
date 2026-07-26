package agenda

import (
	agendaHandler "github.com/ahmadammarm/amami/backend/internal/handler/agenda"
	agendaRepo "github.com/ahmadammarm/amami/backend/internal/repository/agenda"
	agendaSvc "github.com/ahmadammarm/amami/backend/internal/service/agenda"
	"gorm.io/gorm"
)

type Module struct {
	Handler *agendaHandler.AgendaHandler
}

func NewModule(db *gorm.DB) *Module {
	repo := agendaRepo.NewAgendaRepository(db)
	svc := agendaSvc.NewAgendaService(repo)
	handler := agendaHandler.NewAgendaHandler(svc)

	return &Module{
		Handler: handler,
	}
}
