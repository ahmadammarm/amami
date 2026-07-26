package finance

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/dto/finance"
	financeRepo "github.com/ahmadammarm/amami/backend/internal/repository/finance"
	"github.com/google/uuid"
)

type FinanceService interface {
	CreateFund(req finance.CreateFundRequest) (*finance.FundResponse, error)
	GetFunds() ([]finance.FundResponse, error)
	CreateTransaction(req finance.CreateTransactionRequest, userID uuid.UUID) (*finance.TransactionResponse, error)
	GetTransactions(page int, limit int) (*finance.PaginatedTransactionResponse, error)
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
		return nil, err
	}

	return &finance.FundResponse{
		ID:             fund.ID,
		Name:           fund.Name,
		Code:           fund.Code,
		CurrentBalance: fund.CurrentBalance,
	}, nil
}

func (s *financeService) GetFunds() ([]finance.FundResponse, error) {
	funds, err := s.repo.GetFunds()
	if err != nil {
		return nil, err
	}

	var responses []finance.FundResponse
	for _, f := range funds {
		responses = append(responses, finance.FundResponse{
			ID:             f.ID,
			Name:           f.Name,
			Code:           f.Code,
			CurrentBalance: f.CurrentBalance,
		})
	}
	return responses, nil
}

func (s *financeService) CreateTransaction(req finance.CreateTransactionRequest, userID uuid.UUID) (*finance.TransactionResponse, error) {
	// Basic validation happens at handler, but we must enforce business logic here
	if req.Type != "CREDIT" && req.Type != "DEBIT" {
		return nil, errors.New("invalid transaction type")
	}

	tx := &domain.Transaction{
		FundID:      req.FundID,
		Type:        req.Type,
		Amount:      req.Amount,
		Category:    req.Category,
		ReferenceID: req.ReferenceID,
		CreatedBy:   userID,
	}

	// For metadata we might want to store the description
	if req.Description != "" {
		tx.Metadata = []byte(`{"description":"` + req.Description + `"}`)
	} else {
		tx.Metadata = []byte(`{}`)
	}

	if err := s.repo.CreateTransaction(tx); err != nil {
		return nil, err
	}

	return &finance.TransactionResponse{
		ID:          tx.ID,
		FundID:      tx.FundID,
		Type:        tx.Type,
		Amount:      tx.Amount,
		Category:    tx.Category,
		ReferenceID: tx.ReferenceID,
		CreatedBy:   tx.CreatedBy,
		CreatedAt:   tx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *financeService) GetTransactions(page int, limit int) (*finance.PaginatedTransactionResponse, error) {
	transactions, total, err := s.repo.GetTransactions(page, limit)
	if err != nil {
		return nil, err
	}

	var responses []finance.TransactionResponse
	for _, tx := range transactions {
		responses = append(responses, finance.TransactionResponse{
			ID:          tx.ID,
			FundID:      tx.FundID,
			Type:        tx.Type,
			Amount:      tx.Amount,
			Category:    tx.Category,
			ReferenceID: tx.ReferenceID,
			CreatedBy:   tx.CreatedBy,
			CreatedAt:   tx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	
	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return &finance.PaginatedTransactionResponse{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
