package jamaah

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/dto/jamaah"
	jamaahRepo "github.com/ahmadammarm/amami/backend/internal/repository/jamaah"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JamaahService interface {
	CreateJamaah(req jamaah.JamaahRequest) (*jamaah.JamaahResponse, error)
	UpdateJamaah(id uuid.UUID, req jamaah.JamaahRequest) (*jamaah.JamaahResponse, error)
	GetJamaahByID(id uuid.UUID) (*jamaah.JamaahResponse, error)
	GetJamaahs(page int, limit int, query jamaah.ListJamaahQuery) (utils.Page[jamaah.JamaahResponse], error)
	GetSurvey(jamaahID uuid.UUID) (*jamaah.MustahikSurveyResponse, error)
	SubmitSurvey(jamaahID uuid.UUID, req jamaah.MustahikSurveyRequest, surveyorID uuid.UUID) (*jamaah.MustahikSurveyResponse, error)
}

var orderClauses = map[jamaah.JamaahSort]string{
	jamaah.SortNewest:    "created_at desc",
	jamaah.SortScoreDesc: "mustahik_score desc, full_name asc",
	jamaah.SortNameAsc:   "full_name asc",
}

type jamaahService struct {
	repo jamaahRepo.JamaahRepository
}

func NewJamaahService(repo jamaahRepo.JamaahRepository) JamaahService {
	return &jamaahService{repo: repo}
}

func (s *jamaahService) CreateJamaah(req jamaah.JamaahRequest) (*jamaah.JamaahResponse, error) {
	j := &domain.Jamaah{
		FullName:   req.FullName,
		Phone:      req.Phone,
		Address:    req.Address,
		IsMustahik: req.IsMustahik,
	}

	if err := s.repo.Create(j); err != nil {
		return nil, apperror.Internal("failed to create jamaah", err)
	}

	res := jamaah.ToJamaahResponse(j)
	return &res, nil
}

func (s *jamaahService) UpdateJamaah(id uuid.UUID, req jamaah.JamaahRequest) (*jamaah.JamaahResponse, error) {
	j, err := s.findJamaah(id)
	if err != nil {
		return nil, err
	}

	j.FullName = req.FullName
	j.Phone = req.Phone
	j.Address = req.Address
	j.IsMustahik = req.IsMustahik

	if err := s.repo.Update(j); err != nil {
		return nil, apperror.Internal("failed to update jamaah", err)
	}

	res := jamaah.ToJamaahResponse(j)
	return &res, nil
}

func (s *jamaahService) GetJamaahByID(id uuid.UUID) (*jamaah.JamaahResponse, error) {
	j, err := s.findJamaah(id)
	if err != nil {
		return nil, err
	}

	res := jamaah.ToJamaahResponse(j)
	return &res, nil
}

func (s *jamaahService) GetJamaahs(page int, limit int, query jamaah.ListJamaahQuery) (utils.Page[jamaah.JamaahResponse], error) {
	orderBy, ok := orderClauses[query.Sort]
	if !ok {
		orderBy = orderClauses[jamaah.SortNewest]
	}

	jamaahs, total, err := s.repo.GetPaginated(page, limit, jamaahRepo.ListFilter{
		Query:    query.Query,
		Mustahik: query.Mustahik,
		MinScore: query.MinScore,
		OrderBy:  orderBy,
	})
	if err != nil {
		return utils.Page[jamaah.JamaahResponse]{}, apperror.Internal("failed to fetch jamaah", err)
	}

	res := make([]jamaah.JamaahResponse, 0, len(jamaahs))
	for i := range jamaahs {
		res = append(res, jamaah.ToJamaahResponse(&jamaahs[i]))
	}
	return utils.NewPage(res, total, page, limit), nil
}

func (s *jamaahService) GetSurvey(jamaahID uuid.UUID) (*jamaah.MustahikSurveyResponse, error) {
	j, err := s.findJamaah(jamaahID)
	if err != nil {
		return nil, err
	}

	survey, err := s.repo.GetSurvey(jamaahID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("no mustahik survey recorded for this jamaah")
		}
		return nil, apperror.Internal("failed to fetch survey", err)
	}

	res := jamaah.ToSurveyResponse(survey, j.MustahikScore)
	return &res, nil
}

// SubmitSurvey stores the survey answers and recalculates the poverty score.
func (s *jamaahService) SubmitSurvey(jamaahID uuid.UUID, req jamaah.MustahikSurveyRequest, surveyorID uuid.UUID) (*jamaah.MustahikSurveyResponse, error) {
	if _, err := s.findJamaah(jamaahID); err != nil {
		return nil, err
	}

	survey := &domain.MustahikData{
		JamaahID:        jamaahID,
		TotalIncome:     req.TotalIncome,
		DependentsCount: req.DependentsCount,
		HouseStatus:     req.HouseStatus,
		HasDisability:   req.HasDisability,
		CriteriaNotes:   req.CriteriaNotes,
		SurveyedBy:      &surveyorID,
	}
	score := CalculateMustahikScore(MustahikInput{
		TotalIncome:     req.TotalIncome,
		DependentsCount: req.DependentsCount,
		HouseStatus:     req.HouseStatus,
		HasDisability:   req.HasDisability,
	})

	if err := s.repo.SaveSurvey(survey, score); err != nil {
		return nil, apperror.Internal("failed to save survey", err)
	}

	res := jamaah.ToSurveyResponse(survey, score)
	return &res, nil
}

func (s *jamaahService) findJamaah(id uuid.UUID) (*domain.Jamaah, error) {
	j, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("jamaah not found")
		}
		return nil, apperror.Internal("failed to fetch jamaah", err)
	}
	return j, nil
}
