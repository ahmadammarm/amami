package finance

import "github.com/google/uuid"

type CreateFundRequest struct {
	Name string `json:"name" binding:"required,max=100"`
	Code string `json:"code" binding:"required,max=20"`
}

type CreateTransactionRequest struct {
	FundID      uint       `json:"fund_id" binding:"required"`
	Type        string     `json:"type" binding:"required,oneof=DEBIT CREDIT"`
	Amount      int64      `json:"amount" binding:"required,gt=0"`
	Category    string     `json:"category" binding:"required,max=30"`
	Description string     `json:"description" binding:"max=255"`
	ReferenceID *uuid.UUID `json:"reference_id,omitempty"`
}

type FundResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	CurrentBalance int64  `json:"current_balance"`
}

type TransactionResponse struct {
	ID          uuid.UUID  `json:"id"`
	FundID      uint       `json:"fund_id"`
	Type        string     `json:"type"`
	Amount      int64      `json:"amount"`
	Category    string     `json:"category"`
	ReferenceID *uuid.UUID `json:"reference_id,omitempty"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	CreatedAt   string     `json:"created_at"`
}

type PaginatedTransactionResponse struct {
	Data       []TransactionResponse `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}
