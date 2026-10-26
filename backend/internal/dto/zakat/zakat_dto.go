package zakat

import (
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
)

type CollectZakatRequest struct {
	MuzakkiID   uuid.UUID        `json:"muzakki_id" binding:"required"`
	ZakatType   string           `json:"zakat_type" binding:"required,oneof=FITRAH_BERAS MAAL INFAQ"`
	FundID      *uint            `json:"fund_id"` // set only for cash donations posted to the ledger
	AmountOrQty float64          `json:"amount" binding:"required,gt=0"`
	Unit        domain.ZakatUnit `json:"unit" binding:"required,oneof=IDR KG LITER"`
	Description string           `json:"description" binding:"max=1000"`
}

type DistributeZakatRequest struct {
	MustahikID  uuid.UUID        `json:"mustahik_id" binding:"required"`
	FundID      *uint            `json:"fund_id"` // set only for cash distributions posted to the ledger
	AmountOrQty float64          `json:"amount" binding:"required,gt=0"`
	Unit        domain.ZakatUnit `json:"unit" binding:"required,oneof=IDR KG LITER"`
	Description string           `json:"description" binding:"max=1000"`
}

type ZakatDonationResponse struct {
	ID            uuid.UUID        `json:"id"`
	TransactionID *uuid.UUID       `json:"transaction_id"`
	MuzakkiID     uuid.UUID        `json:"muzakki_id"`
	MuzakkiName   string           `json:"muzakki_name"`
	ZakatType     string           `json:"zakat_type"`
	Amount        float64          `json:"amount"`
	Unit          domain.ZakatUnit `json:"unit"`
	Description   string           `json:"description"`
	CreatedAt     time.Time        `json:"created_at"`
}

type ZakatDistributionResponse struct {
	ID            uuid.UUID        `json:"id"`
	TransactionID *uuid.UUID       `json:"transaction_id"`
	MustahikID    uuid.UUID        `json:"mustahik_id"`
	MustahikName  string           `json:"mustahik_name"`
	Amount        float64          `json:"amount"`
	Unit          domain.ZakatUnit `json:"unit"`
	Description   string           `json:"description"`
	DistributedAt time.Time        `json:"distributed_at"`
}

func ToDonationResponse(d *domain.ZakatDonation) ZakatDonationResponse {
	return ZakatDonationResponse{
		ID:            d.ID,
		TransactionID: d.TransactionID,
		MuzakkiID:     d.MuzakkiID,
		MuzakkiName:   d.Muzakki.FullName,
		ZakatType:     d.ZakatType,
		Amount:        d.AmountOrQty,
		Unit:          d.Unit,
		Description:   d.Description,
		CreatedAt:     d.CreatedAt,
	}
}

func ToDistributionResponse(d *domain.ZakatDistribution) ZakatDistributionResponse {
	return ZakatDistributionResponse{
		ID:            d.ID,
		TransactionID: d.TransactionID,
		MustahikID:    d.MustahikID,
		MustahikName:  d.Mustahik.FullName,
		Amount:        d.AmountOrQty,
		Unit:          d.Unit,
		Description:   d.Description,
		DistributedAt: d.DistributedAt,
	}
}
