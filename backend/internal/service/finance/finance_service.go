package finance

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/dto/finance"
	financeRepo "github.com/ahmadammarm/amami/backend/internal/repository/finance"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FinanceService interface {
	CreateFund(req finance.CreateFundRequest) (*finance.FundResponse, error)
	GetFunds() ([]finance.FundResponse, error)
	CreateTransaction(req finance.CreateTransactionRequest, userID uuid.UUID) (*finance.TransactionResponse, error)
	ReverseTransaction(id uuid.UUID, req finance.ReverseTransactionRequest, userID uuid.UUID) (*finance.TransactionResponse, error)
	GetTransactions(page int, limit int) (utils.Page[finance.TransactionResponse], error)
}

type financeService struct {
	repo financeRepo.FinanceRepository
}

func NewFinanceService(repo financeRepo.FinanceRepository) FinanceService {
	return &financeService{repo: repo}
}

func (s *financeService) CreateFund(req finance.CreateFundRequest) (*finance.FundResponse, error) {
	fund := &domain.Fund{
		Name: req.Name,
		Code: req.Code,
	}

	if err := s.repo.CreateFund(fund); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("fund code already exists")
		}
		return nil, apperror.Internal("failed to create fund", err)
	}

	res := finance.ToFundResponse(fund)
	return &res, nil
}

func (s *financeService) GetFunds() ([]finance.FundResponse, error) {
	funds, err := s.repo.GetFunds()
	if err != nil {
		return nil, apperror.Internal("failed to fetch funds", err)
	}

	res := make([]finance.FundResponse, 0, len(funds))
	for i := range funds {
		res = append(res, finance.ToFundResponse(&funds[i]))
	}
	return res, nil
}

func (s *financeService) CreateTransaction(req finance.CreateTransactionRequest, userID uuid.UUID) (*finance.TransactionResponse, error) {
	tx := &domain.Transaction{
		FundID:      req.FundID,
		Type:        req.Type,
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
		ReferenceID: req.ReferenceID,
		CreatedBy:   userID,
	}

	if err := s.repo.CreateTransaction(tx); err != nil {
		return nil, MapLedgerError(err)
	}

	res := finance.ToTransactionResponse(tx, nil)
	return &res, nil
}

// ReverseTransaction posts an offsetting entry for an existing transaction.
// The ledger is append-only, so this is the only way to correct a mistake.
func (s *financeService) ReverseTransaction(id uuid.UUID, req finance.ReverseTransactionRequest, userID uuid.UUID) (*finance.TransactionResponse, error) {
	reversal, err := s.repo.ReverseTransaction(id, func(original *domain.Transaction) (*domain.Transaction, error) {
		return BuildReversal(original, req.Reason, userID)
	})
	if err != nil {
		return nil, MapLedgerError(err)
	}

	res := finance.ToTransactionResponse(reversal, nil)
	return &res, nil
}

// BuildReversal returns the entry that cancels original: same fund and amount,
// opposite direction, linked through ReversalOfID.
func BuildReversal(original *domain.Transaction, reason string, userID uuid.UUID) (*domain.Transaction, error) {
	if original.ReversalOfID != nil {
		return nil, domain.ErrCannotReverseReversal
	}

	originalID := original.ID
	return &domain.Transaction{
		FundID:       original.FundID,
		Type:         original.Type.Opposite(),
		Amount:       original.Amount,
		Category:     domain.CategoryReversal,
		Description:  reason,
		ReversalOfID: &originalID,
		CreatedBy:    userID,
	}, nil
}

func (s *financeService) GetTransactions(page int, limit int) (utils.Page[finance.TransactionResponse], error) {
	rows, total, err := s.repo.GetTransactions(page, limit)
	if err != nil {
		return utils.Page[finance.TransactionResponse]{}, apperror.Internal("failed to fetch transactions", err)
	}

	res := make([]finance.TransactionResponse, 0, len(rows))
	for i := range rows {
		res = append(res, finance.ToTransactionResponse(&rows[i].Transaction, rows[i].ReversedByID))
	}
	return utils.NewPage(res, total, page, limit), nil
}

// MapLedgerError converts ledger rule violations into client-facing errors.
// It is shared with other modules that post to the ledger (e.g. zakat).
func MapLedgerError(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return apperror.NotFound("transaction not found")
	case errors.Is(err, domain.ErrFundNotFound):
		return apperror.Validation(domain.ErrFundNotFound.Error())
	case errors.Is(err, domain.ErrInsufficientBalance),
		errors.Is(err, domain.ErrAlreadyReversed),
		errors.Is(err, domain.ErrCannotReverseReversal):
		return apperror.Conflict(err.Error())
	default:
		return apperror.Internal("failed to post ledger transaction", err)
	}
}
