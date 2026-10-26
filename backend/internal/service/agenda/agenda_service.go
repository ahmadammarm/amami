package agenda

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	agendaDto "github.com/ahmadammarm/amami/backend/internal/dto/agenda"
	agendaRepo "github.com/ahmadammarm/amami/backend/internal/repository/agenda"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/ahmadammarm/amami/backend/pkg/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgendaService interface {
	CreateAgenda(req agendaDto.AgendaRequest, userID uuid.UUID) (*agendaDto.AgendaResponse, error)
	GetAgendas() ([]agendaDto.AgendaResponse, error)
	UpdateAgenda(id uuid.UUID, req agendaDto.AgendaRequest) (*agendaDto.AgendaResponse, error)
	DeleteAgenda(id uuid.UUID) error

	GetDocumentation(agendaID uuid.UUID) ([]agendaDto.DocumentationResponse, error)
	UploadDocumentation(agendaID uuid.UUID, files []UploadFile, userID uuid.UUID) ([]agendaDto.DocumentationResponse, error)
	DeleteDocumentation(agendaID, docID uuid.UUID) error
}

type agendaService struct {
	repo           agendaRepo.AgendaRepository
	storage        storage.Storage
	maxUploadBytes int64
}

func NewAgendaService(repo agendaRepo.AgendaRepository, store storage.Storage, maxUploadBytes int64) AgendaService {
	return &agendaService{repo: repo, storage: store, maxUploadBytes: maxUploadBytes}
}

func (s *agendaService) CreateAgenda(req agendaDto.AgendaRequest, userID uuid.UUID) (*agendaDto.AgendaResponse, error) {
	agenda := &domain.Agenda{CreatedByID: &userID}
	applyRequest(agenda, req)

	if err := s.repo.CreateAgenda(agenda); err != nil {
		return nil, apperror.Internal("failed to create agenda", err)
	}

	res := agendaDto.ToAgendaResponse(agenda)
	return &res, nil
}

func (s *agendaService) GetAgendas() ([]agendaDto.AgendaResponse, error) {
	agendas, err := s.repo.GetAgendas()
	if err != nil {
		return nil, apperror.Internal("failed to fetch agendas", err)
	}

	res := make([]agendaDto.AgendaResponse, 0, len(agendas))
	for i := range agendas {
		res = append(res, agendaDto.ToAgendaResponse(&agendas[i]))
	}
	return res, nil
}

func (s *agendaService) UpdateAgenda(id uuid.UUID, req agendaDto.AgendaRequest) (*agendaDto.AgendaResponse, error) {
	agenda, err := s.findAgenda(id)
	if err != nil {
		return nil, err
	}

	applyRequest(agenda, req)
	if err := s.repo.UpdateAgenda(agenda); err != nil {
		return nil, apperror.FromStore(err, "agenda not found", "failed to update agenda")
	}

	res := agendaDto.ToAgendaResponse(agenda)
	return &res, nil
}

func (s *agendaService) DeleteAgenda(id uuid.UUID) error {
	if err := s.repo.DeleteAgenda(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("agenda not found")
		}
		return apperror.Internal("failed to delete agenda", err)
	}
	return nil
}

func (s *agendaService) findAgenda(id uuid.UUID) (*domain.Agenda, error) {
	agenda, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("agenda not found")
		}
		return nil, apperror.Internal("failed to fetch agenda", err)
	}
	return agenda, nil
}

func applyRequest(agenda *domain.Agenda, req agendaDto.AgendaRequest) {
	agenda.Title = req.Title
	agenda.Description = req.Description
	agenda.Day = req.Day
	agenda.Time = req.Time
	agenda.Location = req.Location
	agenda.Status = req.Status
}
