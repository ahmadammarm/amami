package finance

import (
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
)

type CreateFundRequest struct {
	Name string `json:"name" binding:"required,max=100"`
	Code string `json:"code" binding:"required,max=20"`
}

type CreateTransactionRequest struct {
	FundID      uint                   `json:"fund_id" binding:"required"`
	Type        domain.TransactionType `json:"type" binding:"required,oneof=DEBIT CREDIT"`
	Amount      int64                  `json:"amount" binding:"required,gt=0"`
	Category    string                 `json:"category" binding:"required,max=30"`
	Description string                 `json:"description" binding:"max=255"`
	ReferenceID *uuid.UUID             `json:"reference_id,omitempty"`
}

type ReverseTransactionRequest struct {
	Reason string `json:"reason" binding:"required,min=3,max=255"`
}

type FundResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	CurrentBalance int64  `json:"current_balance"`
}

type TransactionResponse struct {
	ID           uuid.UUID              `json:"id"`
	FundID       uint                   `json:"fund_id"`
	Type         domain.TransactionType `json:"type"`
	Amount       int64                  `json:"amount"`
	Category     string                 `json:"category"`
	Description  string                 `json:"description"`
	ReferenceID  *uuid.UUID             `json:"reference_id"`
	ReversalOfID *uuid.UUID             `json:"reversal_of_id"`
	ReversedByID *uuid.UUID             `json:"reversed_by_id"`
	CreatedBy    uuid.UUID              `json:"created_by"`
	CreatedAt    time.Time              `json:"created_at"`
}

func ToFundResponse(f *domain.Fund) FundResponse {
	return FundResponse{
		ID:             f.ID,
		Name:           f.Name,
		Code:           f.Code,
		CurrentBalance: f.CurrentBalance,
	}
}

// ToTransactionResponse maps a ledger row; reversedByID is the reversal entry
// pointing at this transaction, if any.
func ToTransactionResponse(tx *domain.Transaction, reversedByID *uuid.UUID) TransactionResponse {
	return TransactionResponse{
		ID:           tx.ID,
		FundID:       tx.FundID,
		Type:         tx.Type,
		Amount:       tx.Amount,
		Category:     tx.Category,
		Description:  tx.Description,
		ReferenceID:  tx.ReferenceID,
		ReversalOfID: tx.ReversalOfID,
		ReversedByID: reversedByID,
		CreatedBy:    tx.CreatedBy,
		CreatedAt:    tx.CreatedAt,
	}
}
