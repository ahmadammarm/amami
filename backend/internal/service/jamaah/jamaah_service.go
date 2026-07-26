package jamaah

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/dto/jamaah"
	jamaahRepo "github.com/ahmadammarm/amami/backend/internal/repository/jamaah"
	"github.com/google/uuid"
)

type JamaahService interface {
	CreateJamaah(req jamaah.CreateJamaahRequest) (*jamaah.JamaahResponse, error)
	UpdateJamaah(id uuid.UUID, req jamaah.UpdateJamaahRequest) (*jamaah.JamaahResponse, error)
	GetJamaahByID(id uuid.UUID) (*jamaah.JamaahResponse, error)
	GetJamaahs(page int, limit int, query string) (*jamaah.PaginatedJamaahResponse, error)
}

type jamaahService struct {
	repo jamaahRepo.JamaahRepository
}

func NewJamaahService(repo jamaahRepo.JamaahRepository) JamaahService {
	return &jamaahService{repo: repo}
}

func (s *jamaahService) mapToResponse(j *domain.Jamaah) jamaah.JamaahResponse {
	return jamaah.JamaahResponse{
		ID:            j.ID,
		UserID:        j.UserID,
		FullName:      j.FullName,
		Phone:         j.Phone,
		Address:       j.Address,
		IsMustahik:    j.IsMustahik,
		MustahikScore: j.MustahikScore,
		CreatedAt:     j.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     j.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *jamaahService) CreateJamaah(req jamaah.CreateJamaahRequest) (*jamaah.JamaahResponse, error) {
	j := &domain.Jamaah{
		FullName:   req.FullName,
		Phone:      req.Phone,
		Address:    req.Address,
		IsMustahik: req.IsMustahik,
	}

	if err := s.repo.Create(j); err != nil {
		return nil, err
	}

	res := s.mapToResponse(j)
	return &res, nil
}

func (s *jamaahService) UpdateJamaah(id uuid.UUID, req jamaah.UpdateJamaahRequest) (*jamaah.JamaahResponse, error) {
	j, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("jamaah not found")
	}

	j.FullName = req.FullName
	j.Phone = req.Phone
	j.Address = req.Address
	j.IsMustahik = req.IsMustahik

	if err := s.repo.Update(j); err != nil {
		return nil, err
	}

	res := s.mapToResponse(j)
	return &res, nil
}

func (s *jamaahService) GetJamaahByID(id uuid.UUID) (*jamaah.JamaahResponse, error) {
	j, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("jamaah not found")
	}

	res := s.mapToResponse(j)
	return &res, nil
}

func (s *jamaahService) GetJamaahs(page int, limit int, query string) (*jamaah.PaginatedJamaahResponse, error) {
	jamaahs, total, err := s.repo.GetPaginated(page, limit, query)
	if err != nil {
		return nil, err
	}

	var responses []jamaah.JamaahResponse
	for _, j := range jamaahs {
		responses = append(responses, s.mapToResponse(&j))
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return &jamaah.PaginatedJamaahResponse{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
