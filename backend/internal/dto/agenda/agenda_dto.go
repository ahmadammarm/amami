package agenda

import (
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
)

// AgendaRequest is used for both create and update; every field is replaced on update.
type AgendaRequest struct {
	Title       string              `json:"title" binding:"required,max=200"`
	Description string              `json:"description" binding:"max=5000"`
	Day         string              `json:"day" binding:"required,max=20"`
	Time        string              `json:"time" binding:"required,max=20"`
	Location    string              `json:"location" binding:"max=100"`
	Status      domain.AgendaStatus `json:"status" binding:"required,oneof=SCHEDULED ONGOING COMPLETED CANCELLED"`
}

type AgendaResponse struct {
	ID          uuid.UUID           `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Day         string              `json:"day"`
	Time        string              `json:"time"`
	Location    string              `json:"location"`
	Status      domain.AgendaStatus `json:"status"`
	CreatedByID *uuid.UUID          `json:"created_by_id"`
}

type DocumentationResponse struct {
	ID         uuid.UUID  `json:"id"`
	AgendaID   uuid.UUID  `json:"agenda_id"`
	FileURL    string     `json:"file_url"`
	FileName   string     `json:"file_name"`
	FileType   string     `json:"file_type"`
	FileSize   int64      `json:"file_size"`
	UploadedBy *uuid.UUID `json:"uploaded_by"`
	UploadedAt time.Time  `json:"uploaded_at"`
}

func ToAgendaResponse(a *domain.Agenda) AgendaResponse {
	return AgendaResponse{
		ID:          a.ID,
		Title:       a.Title,
		Description: a.Description,
		Day:         a.Day,
		Time:        a.Time,
		Location:    a.Location,
		Status:      a.Status,
		CreatedByID: a.CreatedByID,
	}
}

func ToDocumentationResponse(d *domain.AgendaDocumentation) DocumentationResponse {
	return DocumentationResponse{
		ID:         d.ID,
		AgendaID:   d.AgendaID,
		FileURL:    d.FileURL,
		FileName:   d.FileName,
		FileType:   d.FileType,
		FileSize:   d.FileSize,
		UploadedBy: d.UploadedBy,
		UploadedAt: d.UploadedAt,
	}
}
