package finance

import (
	"github.com/ahmadammarm/amami/backend/internal/dto/finance"
	financeSvc "github.com/ahmadammarm/amami/backend/internal/service/finance"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type FinanceHandler struct {
	svc financeSvc.FinanceService
}

func NewFinanceHandler(svc financeSvc.FinanceService) *FinanceHandler {
	return &FinanceHandler{svc: svc}
}

func (h *FinanceHandler) CreateFund(c *gin.Context) {
	var req finance.CreateFundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.CreateFund(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.CreatedResponse(c, "Fund created successfully", res)
}

func (h *FinanceHandler) GetFunds(c *gin.Context) {
	res, err := h.svc.GetFunds()
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Funds fetched successfully", res)
}

func (h *FinanceHandler) CreateTransaction(c *gin.Context) {
	userID, err := utils.CurrentUserID(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var req finance.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.CreateTransaction(req, userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.CreatedResponse(c, "Transaction recorded successfully", res)
}

func (h *FinanceHandler) ReverseTransaction(c *gin.Context) {
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

	var req finance.ReverseTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	res, err := h.svc.ReverseTransaction(id, req, userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.CreatedResponse(c, "Transaction reversed successfully", res)
}

func (h *FinanceHandler) GetTransactions(c *gin.Context) {
	page, limit := utils.ParsePagination(c)

	res, err := h.svc.GetTransactions(page, limit)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.SuccessResponse(c, "Transactions fetched successfully", res)
}
