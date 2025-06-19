package handler

import (
	"net/http"
	"strconv"

	"github.com/cihan-ozata/backend-path/internal/domain"
	"github.com/cihan-ozata/backend-path/internal/service"
	"github.com/gin-gonic/gin"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	transactionService *service.TransactionService
}

// NewTransactionHandler creates a new TransactionHandler instance
func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

// CreateTransactionRequest represents the request body for creating a transaction
type CreateTransactionRequest struct {
	UserID      string  `json:"user_id" binding:"required"`
	Type        string  `json:"type" binding:"required,oneof=deposit withdrawal"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

// CreateTransaction handles creating a new transaction
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.ParseUint(req.UserID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	uintUserID := uint(userID)
	transaction := &domain.Transaction{
		FromUserID:  &uintUserID,
		Type:        domain.TransactionType(req.Type),
		Amount:      req.Amount,
		Description: req.Description,
	}

	if err := h.transactionService.CreateTransaction(c.Request.Context(), transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully"})
}

// GetTransaction handles getting a transaction by ID
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	transactionID := c.Param("id")

	transaction, err := h.transactionService.GetTransactionByID(c.Request.Context(), transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if transaction == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// GetUserTransactions handles getting all transactions for a user
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	transactions, err := h.transactionService.GetUserTransactions(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// ProcessTransaction handles processing a transaction
func (h *TransactionHandler) ProcessTransaction(c *gin.Context) {
	transactionID := c.Param("id")

	transaction, err := h.transactionService.GetTransactionByID(c.Request.Context(), transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if transaction == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if err := h.transactionService.ProcessTransaction(c.Request.Context(), transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction processed successfully"})
}

// GetAllTransactions handles getting all transactions (admin only)
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	transactions, err := h.transactionService.GetAllTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

// ListTransactions handles getting all transactions
func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	transactions, err := h.transactionService.GetAllTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
} 