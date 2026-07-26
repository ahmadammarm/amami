package logistics

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
	CurrentStatus string     `json:"current_status" binding:"required"` // GOOD, REPAIR, BROKEN
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

type AgendaResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	Location    string     `json:"location"`
	Status      string     `json:"status"`
	CreatedByID *uuid.UUID `json:"created_by_id"`
}

type CreateAgendaRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	Location    string    `json:"location"`
	Status      string    `json:"status" binding:"required"` // SCHEDULED, ONGOING, COMPLETED, CANCELLED
}
