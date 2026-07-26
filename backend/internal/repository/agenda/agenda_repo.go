package agenda

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"gorm.io/gorm"
)

type AgendaRepository interface {
	CreateAgenda(agenda *domain.Agenda) error
	GetAgendas() ([]domain.Agenda, error)
	UpdateAgenda(agenda *domain.Agenda) error
}

type agendaRepository struct {
	db *gorm.DB
}

func NewAgendaRepository(db *gorm.DB) AgendaRepository {
	return &agendaRepository{db}
}

func (r *agendaRepository) CreateAgenda(agenda *domain.Agenda) error {
	return r.db.Create(agenda).Error
}

func (r *agendaRepository) GetAgendas() ([]domain.Agenda, error) {
	var agendas []domain.Agenda
	err := r.db.Find(&agendas).Error
	return agendas, err
}

func (r *agendaRepository) UpdateAgenda(agenda *domain.Agenda) error {
	return r.db.Save(agenda).Error
}
