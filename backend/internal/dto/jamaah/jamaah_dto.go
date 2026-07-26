package jamaah

import (
	"github.com/google/uuid"
)

type CreateJamaahRequest struct {
	FullName   string  `json:"full_name" binding:"required"`
	Phone      string  `json:"phone"`
	Address    string  `json:"address"`
	IsMustahik bool    `json:"is_mustahik"`
}

type UpdateJamaahRequest struct {
	FullName   string  `json:"full_name" binding:"required"`
	Phone      string  `json:"phone"`
	Address    string  `json:"address"`
	IsMustahik bool    `json:"is_mustahik"`
}

type JamaahResponse struct {
	ID            uuid.UUID `json:"id"`
	UserID        *uuid.UUID `json:"user_id,omitempty"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	IsMustahik    bool      `json:"is_mustahik"`
	MustahikScore float64   `json:"mustahik_score"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

type PaginatedJamaahResponse struct {
	Data       []JamaahResponse `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}
