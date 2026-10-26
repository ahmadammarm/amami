package agenda

import (
	"io"
	"mime/multipart"

	agendaDto "github.com/ahmadammarm/amami/backend/internal/dto/agenda"
	agendaSvc "github.com/ahmadammarm/amami/backend/internal/service/agenda"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

const documentationFormField = "files"

type AgendaHandler struct {
	svc agendaSvc.AgendaService
}

func NewAgendaHandler(svc agendaSvc.AgendaService) *AgendaHandler {
	return &AgendaHandler{svc}
}

func (h *AgendaHandler) CreateAgenda(c *gin.Context) {
	userID, err := utils.CurrentUserID(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req agendaDto.AgendaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.CreateAgenda(req, userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.CreatedResponse(c, "Agenda created successfully", res)
}

func (h *AgendaHandler) GetAgendas(c *gin.Context) {
	res, err := h.svc.GetAgendas()
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Agendas fetched successfully", res)
}

func (h *AgendaHandler) UpdateAgenda(c *gin.Context) {
	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req agendaDto.AgendaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.UpdateAgenda(id, req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Agenda updated successfully", res)
}

func (h *AgendaHandler) DeleteAgenda(c *gin.Context) {
	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err := h.svc.DeleteAgenda(id); err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Agenda deleted successfully", nil)
}

func (h *AgendaHandler) GetDocumentation(c *gin.Context) {
	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	res, err := h.svc.GetDocumentation(id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Documentation fetched successfully", res)
}

func (h *AgendaHandler) UploadDocumentation(c *gin.Context) {
	userID, err := utils.CurrentUserID(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		utils.HandleError(c, apperror.Validation("request must be multipart/form-data"))
		return
	}

	headers := form.File[documentationFormField]
	files := make([]agendaSvc.UploadFile, 0, len(headers))
	for _, fh := range headers {
		files = append(files, toUploadFile(fh))
	}

	res, err := h.svc.UploadDocumentation(id, files, userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.CreatedResponse(c, "Documentation uploaded successfully", res)
}

func (h *AgendaHandler) DeleteDocumentation(c *gin.Context) {
	id, err := utils.ParseUUIDParam(c, "id")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	docID, err := utils.ParseUUIDParam(c, "docId")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err := h.svc.DeleteDocumentation(id, docID); err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Documentation deleted successfully", nil)
}

func toUploadFile(fh *multipart.FileHeader) agendaSvc.UploadFile {
	return agendaSvc.UploadFile{
		Name: fh.Filename,
		Size: fh.Size,
		Open: func() (io.ReadCloser, error) { return fh.Open() },
	}
}
