package logistics

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	logisticsDto "github.com/ahmadammarm/amami/backend/internal/dto/logistics"
	logisticsRepo "github.com/ahmadammarm/amami/backend/internal/repository/logistics"
	"github.com/google/uuid"
)

type LogisticsService interface {
	CreateAsset(req logisticsDto.CreateAssetRequest) (*logisticsDto.AssetResponse, error)
	GetAllAssets(page, limit int) ([]domain.Asset, int64, error)
	UpdateAsset(id string, req logisticsDto.CreateAssetRequest) (*logisticsDto.AssetResponse, error)
	DeleteAsset(id string) error

	CreateLoan(assetID uuid.UUID, req logisticsDto.CreateAssetLoanRequest) error
	ReturnLoan(loanID uuid.UUID, req logisticsDto.ReturnAssetLoanRequest) error

	CreateAgenda(req logisticsDto.CreateAgendaRequest, userID uuid.UUID) (*logisticsDto.AgendaResponse, error)
	GetAgendas() ([]logisticsDto.AgendaResponse, error)
	UpdateAgenda(id string, req logisticsDto.CreateAgendaRequest) (*logisticsDto.AgendaResponse, error)
}

type logisticsService struct {
	repo logisticsRepo.LogisticsRepository
}

func NewLogisticsService(repo logisticsRepo.LogisticsRepository) LogisticsService {
	return &logisticsService{repo}
}

func (s *logisticsService) CreateAsset(req logisticsDto.CreateAssetRequest) (*logisticsDto.AssetResponse, error) {
	asset := &domain.Asset{
		Name:          req.Name,
		SKU:           req.SKU,
		PurchaseDate:  req.PurchaseDate,
		PurchasePrice: req.PurchasePrice,
		CurrentStatus: req.CurrentStatus,
		Location:      req.Location,
	}

	if err := s.repo.CreateAsset(asset); err != nil {
		return nil, err
	}

	return &logisticsDto.AssetResponse{
		ID:            asset.ID,
		Name:          asset.Name,
		SKU:           asset.SKU,
		PurchaseDate:  asset.PurchaseDate,
		PurchasePrice: asset.PurchasePrice,
		CurrentStatus: asset.CurrentStatus,
		Location:      asset.Location,
	}, nil
}

func (s *logisticsService) GetAllAssets(page, limit int) ([]domain.Asset, int64, error) {
	return s.repo.GetAllAssets(page, limit)
}

func (s *logisticsService) UpdateAsset(id string, req logisticsDto.CreateAssetRequest) (*logisticsDto.AssetResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid asset id")
	}
	asset, err := s.repo.GetAssetByID(uid)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	asset.Name = req.Name
	asset.SKU = req.SKU
	asset.PurchaseDate = req.PurchaseDate
	asset.PurchasePrice = req.PurchasePrice
	asset.CurrentStatus = req.CurrentStatus
	asset.Location = req.Location

	if err := s.repo.UpdateAsset(asset); err != nil {
		return nil, err
	}

	return &logisticsDto.AssetResponse{
		ID:            asset.ID,
		Name:          asset.Name,
		SKU:           asset.SKU,
		PurchaseDate:  asset.PurchaseDate,
		PurchasePrice: asset.PurchasePrice,
		CurrentStatus: asset.CurrentStatus,
		Location:      asset.Location,
	}, nil
}

func (s *logisticsService) DeleteAsset(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteAsset(uid)
}

func (s *logisticsService) CreateLoan(assetID uuid.UUID, req logisticsDto.CreateAssetLoanRequest) error {
	asset, err := s.repo.GetAssetByID(assetID)
	if err != nil {
		return err
	}
	if asset.CurrentStatus == "REPAIR" || asset.CurrentStatus == "BROKEN" {
		return errors.New("cannot loan asset that is in repair or broken")
	}

	loan := &domain.AssetLoan{
		AssetID:        assetID,
		JamaahID:       req.JamaahID,
		LoanDate:       req.LoanDate,
		DueDate:        req.DueDate,
		ConditionNotes: req.ConditionNotes,
	}
	return s.repo.CreateLoan(loan)
}

func (s *logisticsService) ReturnLoan(loanID uuid.UUID, req logisticsDto.ReturnAssetLoanRequest) error {
	loan, err := s.repo.GetLoanByID(loanID)
	if err != nil {
		return err
	}
	loan.ReturnDate = &req.ReturnDate
	if req.ConditionNotes != "" {
		loan.ConditionNotes = loan.ConditionNotes + " | Return: " + req.ConditionNotes
	}
	return s.repo.UpdateLoan(loan)
}

func (s *logisticsService) CreateAgenda(req logisticsDto.CreateAgendaRequest, userID uuid.UUID) (*logisticsDto.AgendaResponse, error) {
	agenda := &domain.Agenda{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Location:    req.Location,
		Status:      req.Status,
		CreatedByID: &userID,
	}
	if err := s.repo.CreateAgenda(agenda); err != nil {
		return nil, err
	}
	return &logisticsDto.AgendaResponse{
		ID:          agenda.ID,
		Title:       agenda.Title,
		Description: agenda.Description,
		StartTime:   agenda.StartTime,
		EndTime:     agenda.EndTime,
		Location:    agenda.Location,
		Status:      agenda.Status,
		CreatedByID: agenda.CreatedByID,
	}, nil
}

func (s *logisticsService) GetAgendas() ([]logisticsDto.AgendaResponse, error) {
	agendas, err := s.repo.GetAgendas()
	if err != nil {
		return nil, err
	}

	res := make([]logisticsDto.AgendaResponse, 0, len(agendas))
	for _, a := range agendas {
		res = append(res, logisticsDto.AgendaResponse{
			ID:          a.ID,
			Title:       a.Title,
			Description: a.Description,
			StartTime:   a.StartTime,
			EndTime:     a.EndTime,
			Location:    a.Location,
			Status:      a.Status,
			CreatedByID: a.CreatedByID,
		})
	}
	return res, nil
}

func (s *logisticsService) UpdateAgenda(id string, req logisticsDto.CreateAgendaRequest) (*logisticsDto.AgendaResponse, error) {
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
	agenda.StartTime = req.StartTime
	agenda.EndTime = req.EndTime
	agenda.Location = req.Location
	agenda.Status = req.Status

	if err := s.repo.UpdateAgenda(agenda); err != nil {
		return nil, err
	}

	return &logisticsDto.AgendaResponse{
		ID:          agenda.ID,
		Title:       agenda.Title,
		Description: agenda.Description,
		StartTime:   agenda.StartTime,
		EndTime:     agenda.EndTime,
		Location:    agenda.Location,
		Status:      agenda.Status,
		CreatedByID: agenda.CreatedByID,
	}, nil
}
