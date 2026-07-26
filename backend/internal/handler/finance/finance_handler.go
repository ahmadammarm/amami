package finance

import (
	"net/http"
	"strconv"

	"github.com/ahmadammarm/amami/backend/internal/dto/finance"
	financeSvc "github.com/ahmadammarm/amami/backend/internal/service/finance"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.CreateFund(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *FinanceHandler) GetFunds(c *gin.Context) {
	res, err := h.svc.GetFunds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *FinanceHandler) CreateTransaction(c *gin.Context) {
	var req finance.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		// Fallback if it is stored as string
		parsed, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			return
		}
		userID = parsed
	}

	res, err := h.svc.CreateTransaction(req, userID)
	if err != nil {
		// Differentiate between bad request (e.g. insufficient funds) and internal error
		if err.Error() == "insufficient fund balance" || err.Error() == "invalid transaction type" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *FinanceHandler) GetTransactions(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	res, err := h.svc.GetTransactions(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
