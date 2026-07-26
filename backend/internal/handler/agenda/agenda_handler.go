package agenda

import (
	agendaDto "github.com/ahmadammarm/amami/backend/internal/dto/agenda"
	agendaSvc "github.com/ahmadammarm/amami/backend/internal/service/agenda"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AgendaHandler struct {
	svc agendaSvc.AgendaService
}

func NewAgendaHandler(svc agendaSvc.AgendaService) *AgendaHandler {
	return &AgendaHandler{svc}
}

func (h *AgendaHandler) CreateAgenda(c *gin.Context) {
	var req agendaDto.CreateAgendaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	res, err := h.svc.CreateAgenda(req, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func (h *AgendaHandler) GetAgendas(c *gin.Context) {
	res, err := h.svc.GetAgendas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *AgendaHandler) UpdateAgenda(c *gin.Context) {
	id := c.Param("id")
	var req agendaDto.CreateAgendaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.UpdateAgenda(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
