package jamaah

import (
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
)

type JamaahRequest struct {
	FullName   string `json:"full_name" binding:"required,max=100"`
	Phone      string `json:"phone" binding:"max=20"`
	Address    string `json:"address" binding:"max=1000"`
	IsMustahik bool   `json:"is_mustahik"`
}

// JamaahSort values map to fixed ORDER BY clauses in the repository.
type JamaahSort string

const (
	SortNewest    JamaahSort = "newest"
	SortScoreDesc JamaahSort = "score_desc"
	SortNameAsc   JamaahSort = "name_asc"
)

type ListJamaahQuery struct {
	Query    string     `form:"query" binding:"max=100"`
	Mustahik *bool      `form:"mustahik"`
	MinScore *float64   `form:"min_score" binding:"omitempty,gte=0,lte=100"`
	Sort     JamaahSort `form:"sort" binding:"omitempty,oneof=newest score_desc name_asc"`
}

type MustahikSurveyRequest struct {
	TotalIncome     int64              `json:"total_income" binding:"gte=0"`
	DependentsCount int                `json:"dependents_count" binding:"gte=0,lte=30"`
	HouseStatus     domain.HouseStatus `json:"house_status" binding:"required,oneof=OWNED FAMILY RENTED NONE"`
	HasDisability   bool               `json:"has_disability"`
	CriteriaNotes   string             `json:"criteria_notes" binding:"max=2000"`
}

type JamaahResponse struct {
	ID            uuid.UUID  `json:"id"`
	UserID        *uuid.UUID `json:"user_id"`
	FullName      string     `json:"full_name"`
	Phone         string     `json:"phone"`
	Address       string     `json:"address"`
	IsMustahik    bool       `json:"is_mustahik"`
	MustahikScore float64    `json:"mustahik_score"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type MustahikSurveyResponse struct {
	JamaahID        uuid.UUID          `json:"jamaah_id"`
	TotalIncome     int64              `json:"total_income"`
	DependentsCount int                `json:"dependents_count"`
	HouseStatus     domain.HouseStatus `json:"house_status"`
	HasDisability   bool               `json:"has_disability"`
	CriteriaNotes   string             `json:"criteria_notes"`
	MustahikScore   float64            `json:"mustahik_score"`
	SurveyedBy      *uuid.UUID         `json:"surveyed_by"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

func ToJamaahResponse(j *domain.Jamaah) JamaahResponse {
	return JamaahResponse{
		ID:            j.ID,
		UserID:        j.UserID,
		FullName:      j.FullName,
		Phone:         j.Phone,
		Address:       j.Address,
		IsMustahik:    j.IsMustahik,
		MustahikScore: j.MustahikScore,
		CreatedAt:     j.CreatedAt,
		UpdatedAt:     j.UpdatedAt,
	}
}

func ToSurveyResponse(d *domain.MustahikData, score float64) MustahikSurveyResponse {
	return MustahikSurveyResponse{
		JamaahID:        d.JamaahID,
		TotalIncome:     d.TotalIncome,
		DependentsCount: d.DependentsCount,
		HouseStatus:     d.HouseStatus,
		HasDisability:   d.HasDisability,
		CriteriaNotes:   d.CriteriaNotes,
		MustahikScore:   score,
		SurveyedBy:      d.SurveyedBy,
		UpdatedAt:       d.UpdatedAt,
	}
}
