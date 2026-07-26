package zakat

import (
	"github.com/google/uuid"
)

type CollectZakatRequest struct {
	MuzakkiID   uuid.UUID `json:"muzakki_id" binding:"required"`
	ZakatType   string    `json:"zakat_type" binding:"required,oneof=FITRAH_BERAS MAAL INFAQ"` // "FITRAH_BERAS", "MAAL", "INFAQ"
	FundID      *uint     `json:"fund_id"` // Nullable for non-cash
	AmountOrQty float64   `json:"amount" binding:"required,gt=0"`
	Unit        string    `json:"unit" binding:"required"` // "IDR", "KG", "LITER"
	Description string    `json:"description"`
}

type DistributeZakatRequest struct {
	MustahikID  uuid.UUID `json:"mustahik_id" binding:"required"`
	FundID      *uint     `json:"fund_id"` // Nullable for non-cash
	AmountOrQty float64   `json:"amount" binding:"required,gt=0"`
	Unit        string    `json:"unit" binding:"required"`
	Description string    `json:"description"`
}

type ZakatDonationResponse struct {
	ID            uuid.UUID  `json:"id"`
	TransactionID *uuid.UUID `json:"transaction_id,omitempty"`
	MuzakkiID     uuid.UUID  `json:"muzakki_id"`
	MuzakkiName   string     `json:"muzakki_name"`
	ZakatType     string     `json:"zakat_type"`
	Amount        float64    `json:"amount"`
	Unit          string     `json:"unit"`
	Description   string     `json:"description"`
	CreatedAt     string     `json:"created_at"`
}

type ZakatDistributionResponse struct {
	ID            uuid.UUID  `json:"id"`
	TransactionID *uuid.UUID `json:"transaction_id,omitempty"`
	MustahikID    uuid.UUID  `json:"mustahik_id"`
	MustahikName  string     `json:"mustahik_name"`
	Amount        float64    `json:"amount"`
	Unit          string     `json:"unit"`
	Description   string     `json:"description"`
	DistributedAt string     `json:"distributed_at"`
}

type PaginatedDonationResponse struct {
	Data       []ZakatDonationResponse `json:"data"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalPages int                     `json:"total_pages"`
}

type PaginatedDistributionResponse struct {
	Data       []ZakatDistributionResponse `json:"data"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	Limit      int                         `json:"limit"`
	TotalPages int                         `json:"total_pages"`
}
