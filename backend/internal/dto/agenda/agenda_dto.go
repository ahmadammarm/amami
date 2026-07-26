package agenda

import (
	"github.com/google/uuid"
)

type AgendaResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Day         string     `json:"day"`
	Time        string     `json:"time"`
	Location    string     `json:"location"`
	Status      string     `json:"status"`
	CreatedByID *uuid.UUID `json:"created_by_id"`
}

type CreateAgendaRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Day         string    `json:"day" binding:"required"`
	Time        string    `json:"time" binding:"required"`
	Location    string    `json:"location"`
	Status      string    `json:"status" binding:"required"` // SCHEDULED, ONGOING, COMPLETED, CANCELLED
}
