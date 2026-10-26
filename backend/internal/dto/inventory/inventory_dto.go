package inventory

import (
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
)

type AssetRequest struct {
	Name          string             `json:"name" binding:"required,max=255"`
	SKU           string             `json:"sku" binding:"required,max=50"`
	PurchaseDate  *time.Time         `json:"purchase_date"`
	PurchasePrice int64              `json:"purchase_price" binding:"gte=0"`
	CurrentStatus domain.AssetStatus `json:"current_status" binding:"required,oneof=GOOD LOANED NEEDS_REPAIR REPAIR BROKEN"`
	Location      string             `json:"location" binding:"max=100"`
}

type CreateAssetLoanRequest struct {
	JamaahID       uuid.UUID  `json:"jamaah_id" binding:"required"`
	LoanDate       time.Time  `json:"loan_date" binding:"required"`
	DueDate        *time.Time `json:"due_date"`
	ConditionNotes string     `json:"condition_notes" binding:"max=1000"`
}

type ReturnAssetLoanRequest struct {
	ReturnDate     time.Time `json:"return_date" binding:"required"`
	ConditionNotes string    `json:"condition_notes" binding:"max=1000"`
}

type AssetResponse struct {
	ID            uuid.UUID          `json:"id"`
	Name          string             `json:"name"`
	SKU           string             `json:"sku"`
	PurchaseDate  *time.Time         `json:"purchase_date"`
	PurchasePrice int64              `json:"purchase_price"`
	CurrentStatus domain.AssetStatus `json:"current_status"`
	Location      string             `json:"location"`
}

type AssetLoanResponse struct {
	ID             uuid.UUID  `json:"id"`
	AssetID        uuid.UUID  `json:"asset_id"`
	JamaahID       uuid.UUID  `json:"jamaah_id"`
	LoanDate       time.Time  `json:"loan_date"`
	DueDate        *time.Time `json:"due_date"`
	ReturnDate     *time.Time `json:"return_date"`
	ConditionNotes string     `json:"condition_notes"`
}

func ToAssetResponse(a *domain.Asset) AssetResponse {
	return AssetResponse{
		ID:            a.ID,
		Name:          a.Name,
		SKU:           a.SKU,
		PurchaseDate:  a.PurchaseDate,
		PurchasePrice: a.PurchasePrice,
		CurrentStatus: a.CurrentStatus,
		Location:      a.Location,
	}
}

func ToLoanResponse(l *domain.AssetLoan) AssetLoanResponse {
	return AssetLoanResponse{
		ID:             l.ID,
		AssetID:        l.AssetID,
		JamaahID:       l.JamaahID,
		LoanDate:       l.LoanDate,
		DueDate:        l.DueDate,
		ReturnDate:     l.ReturnDate,
		ConditionNotes: l.ConditionNotes,
	}
}
