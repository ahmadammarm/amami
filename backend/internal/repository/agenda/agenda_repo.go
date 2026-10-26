package agenda

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AgendaRepository interface {
	CreateAgenda(agenda *domain.Agenda) error
	GetAgendas() ([]domain.Agenda, error)
	GetByID(id uuid.UUID) (*domain.Agenda, error)
	UpdateAgenda(agenda *domain.Agenda) error
	DeleteAgenda(id uuid.UUID) error

	GetDocumentation(agendaID uuid.UUID) ([]domain.AgendaDocumentation, error)
	GetDocumentationByID(agendaID, docID uuid.UUID) (*domain.AgendaDocumentation, error)
	CreateDocumentation(docs []domain.AgendaDocumentation) error
	DeleteDocumentation(id uuid.UUID) error
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
	err := r.db.Order("title").Find(&agendas).Error
	return agendas, err
}

func (r *agendaRepository) GetByID(id uuid.UUID) (*domain.Agenda, error) {
	var agenda domain.Agenda
	if err := r.db.First(&agenda, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &agenda, nil
}

// UpdateAgenda reports gorm.ErrRecordNotFound if the agenda was deleted after
// it was read.
func (r *agendaRepository) UpdateAgenda(agenda *domain.Agenda) error {
	return database.ExpectRows(r.db.Model(agenda).
		Select("title", "description", "day", "time", "location", "status").
		Updates(agenda))
}

// DeleteAgenda soft-deletes the agenda; its documentation is kept for the archive.
func (r *agendaRepository) DeleteAgenda(id uuid.UUID) error {
	return database.ExpectRows(r.db.Delete(&domain.Agenda{}, "id = ?", id))
}

func (r *agendaRepository) GetDocumentation(agendaID uuid.UUID) ([]domain.AgendaDocumentation, error) {
	var docs []domain.AgendaDocumentation
	err := r.db.Where("agenda_id = ?", agendaID).Order("uploaded_at desc").Find(&docs).Error
	return docs, err
}

func (r *agendaRepository) GetDocumentationByID(agendaID, docID uuid.UUID) (*domain.AgendaDocumentation, error) {
	var doc domain.AgendaDocumentation
	if err := r.db.First(&doc, "id = ? AND agenda_id = ?", docID, agendaID).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *agendaRepository) CreateDocumentation(docs []domain.AgendaDocumentation) error {
	return r.db.Omit(clause.Associations).Create(&docs).Error
}

func (r *agendaRepository) DeleteDocumentation(id uuid.UUID) error {
	return r.db.Delete(&domain.AgendaDocumentation{}, "id = ?", id).Error
}
