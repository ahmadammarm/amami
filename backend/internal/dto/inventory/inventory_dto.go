package inventory

import (
	"time"
	"github.com/google/uuid"
)

type AssetResponse struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	SKU           string     `json:"sku"`
	PurchaseDate  *time.Time `json:"purchase_date"`
	PurchasePrice int64      `json:"purchase_price"`
	CurrentStatus string     `json:"current_status"`
	Location      string     `json:"location"`
}

type CreateAssetRequest struct {
	Name          string     `json:"name" binding:"required"`
	SKU           string     `json:"sku" binding:"required"`
	PurchaseDate  *time.Time `json:"purchase_date"`
	PurchasePrice int64      `json:"purchase_price"`
	CurrentStatus string     `json:"current_status" binding:"required"` // GOOD, REPAIR, BROKEN, LOANED, NEEDS_REPAIR
	Location      string     `json:"location"`
}

type CreateAssetLoanRequest struct {
	JamaahID       uuid.UUID  `json:"jamaah_id" binding:"required"`
	LoanDate       time.Time  `json:"loan_date" binding:"required"`
	DueDate        *time.Time `json:"due_date"`
	ConditionNotes string     `json:"condition_notes"`
}

type ReturnAssetLoanRequest struct {
	ReturnDate     time.Time `json:"return_date" binding:"required"`
	ConditionNotes string    `json:"condition_notes"`
}
