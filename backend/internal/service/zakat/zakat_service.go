package zakat

import (
	"errors"
	"math"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/dto/zakat"
	zakatRepo "github.com/ahmadammarm/amami/backend/internal/repository/zakat"
	financeSvc "github.com/ahmadammarm/amami/backend/internal/service/finance"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ZakatService interface {
	CollectZakat(req zakat.CollectZakatRequest, userID uuid.UUID) (*zakat.ZakatDonationResponse, error)
	DistributeZakat(req zakat.DistributeZakatRequest, userID uuid.UUID) (*zakat.ZakatDistributionResponse, error)
	GetDonations(page int, limit int) (utils.Page[zakat.ZakatDonationResponse], error)
	GetDistributions(page int, limit int) (utils.Page[zakat.ZakatDistributionResponse], error)
}

type zakatService struct {
	repo zakatRepo.ZakatRepository
}

func NewZakatService(repo zakatRepo.ZakatRepository) ZakatService {
	return &zakatService{repo: repo}
}

func (s *zakatService) CollectZakat(req zakat.CollectZakatRequest, userID uuid.UUID) (*zakat.ZakatDonationResponse, error) {
	txModel, err := buildLedgerEntry(req.FundID, req.AmountOrQty, req.Unit, domain.TransactionCredit, domain.CategoryZakatIn, userID)
	if err != nil {
		return nil, err
	}

	donation := &domain.ZakatDonation{
		MuzakkiID:   req.MuzakkiID,
		ZakatType:   req.ZakatType,
		AmountOrQty: req.AmountOrQty,
		Unit:        req.Unit,
		Description: req.Description,
	}

	if err := s.repo.CollectZakat(donation, txModel); err != nil {
		return nil, financeSvc.MapLedgerError(err)
	}

	res := zakat.ToDonationResponse(donation)
	return &res, nil
}

func (s *zakatService) DistributeZakat(req zakat.DistributeZakatRequest, userID uuid.UUID) (*zakat.ZakatDistributionResponse, error) {
	mustahik, err := s.repo.FindJamaah(req.MustahikID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Validation("recipient not found")
		}
		return nil, apperror.Internal("failed to load recipient", err)
	}
	if !mustahik.IsMustahik {
		return nil, apperror.Validation("recipient is not registered as mustahik")
	}

	txModel, err := buildLedgerEntry(req.FundID, req.AmountOrQty, req.Unit, domain.TransactionDebit, domain.CategoryZakatOut, userID)
	if err != nil {
		return nil, err
	}

	distribution := &domain.ZakatDistribution{
		MustahikID:  req.MustahikID,
		AmountOrQty: req.AmountOrQty,
		Unit:        req.Unit,
		Description: req.Description,
	}

	if err := s.repo.DistributeZakat(distribution, txModel); err != nil {
		return nil, financeSvc.MapLedgerError(err)
	}

	res := zakat.ToDistributionResponse(distribution)
	return &res, nil
}

// buildLedgerEntry returns the ledger transaction for a zakat record, or nil
// when no fund is given (goods such as rice are not posted to the ledger).
// Only whole rupiah amounts can be posted, since the ledger stores int64 IDR.
func buildLedgerEntry(fundID *uint, amount float64, unit domain.ZakatUnit, txType domain.TransactionType, category string, userID uuid.UUID) (*domain.Transaction, error) {
	if fundID == nil {
		return nil, nil
	}
	if unit != domain.ZakatUnitIDR {
		return nil, apperror.Validation("only IDR amounts can be posted to a fund")
	}
	if amount != math.Trunc(amount) || amount > math.MaxInt64 {
		return nil, apperror.Validation("amount posted to a fund must be a whole rupiah value")
	}

	return &domain.Transaction{
		FundID:    *fundID,
		Type:      txType,
		Amount:    int64(amount),
		Category:  category,
		CreatedBy: userID,
	}, nil
}

func (s *zakatService) GetDonations(page int, limit int) (utils.Page[zakat.ZakatDonationResponse], error) {
	donations, total, err := s.repo.GetDonations(page, limit)
	if err != nil {
		return utils.Page[zakat.ZakatDonationResponse]{}, apperror.Internal("failed to fetch donations", err)
	}

	res := make([]zakat.ZakatDonationResponse, 0, len(donations))
	for i := range donations {
		res = append(res, zakat.ToDonationResponse(&donations[i]))
	}
	return utils.NewPage(res, total, page, limit), nil
}

func (s *zakatService) GetDistributions(page int, limit int) (utils.Page[zakat.ZakatDistributionResponse], error) {
	distributions, total, err := s.repo.GetDistributions(page, limit)
	if err != nil {
		return utils.Page[zakat.ZakatDistributionResponse]{}, apperror.Internal("failed to fetch distributions", err)
	}

	res := make([]zakat.ZakatDistributionResponse, 0, len(distributions))
	for i := range distributions {
		res = append(res, zakat.ToDistributionResponse(&distributions[i]))
	}
	return utils.NewPage(res, total, page, limit), nil
}
