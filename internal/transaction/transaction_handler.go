package transaction

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service *transactionService
}

func NewTransactionHandler(service *transactionService) *transactionHandler {
	return &transactionHandler{service: service}
}

type TransactionRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (h *transactionHandler) ProcessTransaction(c *gin.Context) {
	// 1. Get userID from JWT context
	userID := c.GetInt("userID")

	// 2. Parse request body
	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Get full user object — needed for saving config
	user, err := h.service.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	// 4. Process transaction
	tx, err := h.service.ProcessTransaction(userID, req.Amount, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "transaction processed successfully",
		"transaction": tx,
	})
}

func (h *transactionHandler) GetMemberContribution(c *gin.Context) {
	// 1. Get userID from JWT context
	userID := c.GetInt("userID")

	// 2. Get account type from query param
	accountType := c.Query("account_type")
	if accountType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_type is required"})
		return
	}

	// 3. Get contribution
	total, err := h.service.GetMemberContribution(userID, accountType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":      userID,
		"account_type": accountType,
		"total":        total,
	})
}
