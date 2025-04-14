package handler

import (
	"net/http"
	"strconv"

	"github.com/cihan-ozata/backend-path/internal/domain"
	"github.com/cihan-ozata/backend-path/internal/service"
	"github.com/gin-gonic/gin"
)

// WalletHandler handles wallet-related HTTP requests
type WalletHandler struct {
	walletService *service.WalletService
}

// NewWalletHandler creates a new WalletHandler instance
func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

// CreateWalletRequest represents the request body for creating a wallet
type CreateWalletRequest struct {
	UserID   string  `json:"user_id" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
	Balance  float64 `json:"balance" binding:"gte=0"`
}

// CreateWallet handles creating a new wallet
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.ParseUint(req.UserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	wallet := &domain.Wallet{
		UserID:   uint(userID),
		Currency: req.Currency,
		Balance:  req.Balance,
	}

	if err := h.walletService.CreateWallet(c.Request.Context(), wallet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Wallet created successfully"})
}

// GetWallet handles getting a wallet by ID
func (h *WalletHandler) GetWallet(c *gin.Context) {
	walletID := c.Param("id")

	wallet, err := h.walletService.GetWalletByID(c.Request.Context(), walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if wallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// GetUserWallet handles getting a wallet by user ID
func (h *WalletHandler) GetUserWallet(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	wallet, err := h.walletService.GetWalletByUserID(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if wallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// GetWalletBalance handles getting the balance of a wallet by user ID
func (h *WalletHandler) GetWalletBalance(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	balance, err := h.walletService.GetWalletBalance(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
} 