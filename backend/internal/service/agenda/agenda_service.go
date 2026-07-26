package agenda

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	agendaDto "github.com/ahmadammarm/amami/backend/internal/dto/agenda"
	agendaRepo "github.com/ahmadammarm/amami/backend/internal/repository/agenda"
	"github.com/google/uuid"
)

type AgendaService interface {
	CreateAgenda(req agendaDto.CreateAgendaRequest, userID uuid.UUID) (*agendaDto.AgendaResponse, error)
	GetAgendas() ([]agendaDto.AgendaResponse, error)
	UpdateAgenda(id string, req agendaDto.CreateAgendaRequest) (*agendaDto.AgendaResponse, error)
}

type agendaService struct {
	repo agendaRepo.AgendaRepository
}

func NewAgendaService(repo agendaRepo.AgendaRepository) AgendaService {
	return &agendaService{repo}
}

func (s *agendaService) CreateAgenda(req agendaDto.CreateAgendaRequest, userID uuid.UUID) (*agendaDto.AgendaResponse, error) {
	agenda := &domain.Agenda{
		Title:       req.Title,
		Description: req.Description,
		Day:         req.Day,
		Time:        req.Time,
		Location:    req.Location,
		Status:      req.Status,
		CreatedByID: &userID,
	}
	if err := s.repo.CreateAgenda(agenda); err != nil {
		return nil, err
	}
	return &agendaDto.AgendaResponse{
		ID:          agenda.ID,
		Title:       agenda.Title,
		Description: agenda.Description,
		Day:         agenda.Day,
		Time:        agenda.Time,
		Location:    agenda.Location,
		Status:      agenda.Status,
		CreatedByID: agenda.CreatedByID,
	}, nil
}

func (s *agendaService) GetAgendas() ([]agendaDto.AgendaResponse, error) {
	agendas, err := s.repo.GetAgendas()
	if err != nil {
		return nil, err
	}

	res := make([]agendaDto.AgendaResponse, 0, len(agendas))
	for _, a := range agendas {
		res = append(res, agendaDto.AgendaResponse{
			ID:          a.ID,
			Title:       a.Title,
			Description: a.Description,
			Day:         a.Day,
			Time:        a.Time,
			Location:    a.Location,
			Status:      a.Status,
			CreatedByID: a.CreatedByID,
		})
	}
	return res, nil
}

func (s *agendaService) UpdateAgenda(id string, req agendaDto.CreateAgendaRequest) (*agendaDto.AgendaResponse, error) {
	agendas, err := s.repo.GetAgendas()
	if err != nil {
		return nil, err
	}
	var agenda *domain.Agenda
	for i := range agendas {
		if agendas[i].ID.String() == id {
			agenda = &agendas[i]
			break
		}
	}
	if agenda == nil {
		return nil, errors.New("agenda not found")
	}

	agenda.Title = req.Title
	agenda.Description = req.Description
	agenda.Day = req.Day
	agenda.Time = req.Time
	agenda.Location = req.Location
	agenda.Status = req.Status

	if err := s.repo.UpdateAgenda(agenda); err != nil {
		return nil, err
	}

	return &agendaDto.AgendaResponse{
		ID:          agenda.ID,
		Title:       agenda.Title,
		Description: agenda.Description,
		Day:         agenda.Day,
		Time:        agenda.Time,
		Location:    agenda.Location,
		Status:      agenda.Status,
		CreatedByID: agenda.CreatedByID,
	}, nil
}
