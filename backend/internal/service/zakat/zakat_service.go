package zakat

import (
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/dto/zakat"
	zakatRepo "github.com/ahmadammarm/amami/backend/internal/repository/zakat"
	"github.com/google/uuid"
)

type ZakatService interface {
	CollectZakat(req zakat.CollectZakatRequest, userID uuid.UUID) (*zakat.ZakatDonationResponse, error)
	DistributeZakat(req zakat.DistributeZakatRequest, userID uuid.UUID) (*zakat.ZakatDistributionResponse, error)
	GetDonations(page int, limit int) (*zakat.PaginatedDonationResponse, error)
	GetDistributions(page int, limit int) (*zakat.PaginatedDistributionResponse, error)
}

type zakatService struct {
	repo zakatRepo.ZakatRepository
}

func NewZakatService(repo zakatRepo.ZakatRepository) ZakatService {
	return &zakatService{repo: repo}
}

func (s *zakatService) CollectZakat(req zakat.CollectZakatRequest, userID uuid.UUID) (*zakat.ZakatDonationResponse, error) {
	var txModel *domain.Transaction
	if req.FundID != nil {
		txModel = &domain.Transaction{
			FundID:    *req.FundID,
			Type:      "CREDIT",
			Amount:    int64(req.AmountOrQty), // Assuming cash is integer
			Category:  "Penerimaan Zakat",
			CreatedBy: userID,
		}
	}

	donation := &domain.ZakatDonation{
		MuzakkiID:   req.MuzakkiID,
		ZakatType:   req.ZakatType,
		AmountOrQty: req.AmountOrQty,
		Unit:        req.Unit,
		Description: req.Description,
	}

	if err := s.repo.CollectZakat(donation, txModel); err != nil {
		return nil, err
	}

	res := &zakat.ZakatDonationResponse{
		ID:          donation.ID,
		MuzakkiID:   donation.MuzakkiID,
		ZakatType:   donation.ZakatType,
		Amount:      donation.AmountOrQty,
		Unit:        donation.Unit,
		Description: donation.Description,
		CreatedAt:   time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}
	if txModel != nil {
		res.TransactionID = &txModel.ID
	}
	return res, nil
}

func (s *zakatService) DistributeZakat(req zakat.DistributeZakatRequest, userID uuid.UUID) (*zakat.ZakatDistributionResponse, error) {
	var txModel *domain.Transaction
	if req.FundID != nil {
		txModel = &domain.Transaction{
			FundID:    *req.FundID,
			Type:      "DEBIT",
			Amount:    int64(req.AmountOrQty), // Assuming cash is integer
			Category:  "Penyaluran Zakat",
			CreatedBy: userID,
		}
	}

	distribution := &domain.ZakatDistribution{
		MustahikID:  req.MustahikID,
		AmountOrQty: req.AmountOrQty,
		Unit:        req.Unit,
		Description: req.Description,
	}

	if err := s.repo.DistributeZakat(distribution, txModel); err != nil {
		return nil, err
	}

	res := &zakat.ZakatDistributionResponse{
		ID:            distribution.ID,
		MustahikID:    distribution.MustahikID,
		Amount:        distribution.AmountOrQty,
		Unit:        distribution.Unit,
		Description: distribution.Description,
		DistributedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}
	if txModel != nil {
		res.TransactionID = &txModel.ID
	}
	return res, nil
}

func (s *zakatService) GetDonations(page int, limit int) (*zakat.PaginatedDonationResponse, error) {
	donations, total, err := s.repo.GetDonations(page, limit)
	if err != nil {
		return nil, err
	}

	var res []zakat.ZakatDonationResponse
	for _, d := range donations {
		createdAt := d.CreatedAt
		if d.Transaction != nil {
			createdAt = d.Transaction.CreatedAt
		}
		res = append(res, zakat.ZakatDonationResponse{
			ID:            d.ID,
			TransactionID: d.TransactionID,
			MuzakkiID:     d.MuzakkiID,
			MuzakkiName:   d.Muzakki.FullName,
			ZakatType:     d.ZakatType,
			Amount:        d.AmountOrQty,
			Unit:          d.Unit,
			Description:   d.Description,
			CreatedAt:     createdAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return &zakat.PaginatedDonationResponse{
		Data:       res,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *zakatService) GetDistributions(page int, limit int) (*zakat.PaginatedDistributionResponse, error) {
	distributions, total, err := s.repo.GetDistributions(page, limit)
	if err != nil {
		return nil, err
	}

	var res []zakat.ZakatDistributionResponse
	for _, d := range distributions {
		res = append(res, zakat.ZakatDistributionResponse{
			ID:            d.ID,
			TransactionID: d.TransactionID,
			MustahikID:    d.MustahikID,
			MustahikName:  d.Mustahik.FullName,
			Amount:        d.AmountOrQty,
			Unit:          d.Unit,
			Description:   d.Description,
			DistributedAt: d.DistributedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return &zakat.PaginatedDistributionResponse{
		Data:       res,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
