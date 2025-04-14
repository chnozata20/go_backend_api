package domain

import (
	"errors"
	"time"
)

// Wallet represents a user's wallet
type Wallet struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;unique" json:"user_id"`
	Balance   float64   `gorm:"type:decimal(20,2);not null;default:0" json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate performs validation on the Wallet struct
func (w *Wallet) Validate() error {
	if w.UserID == 0 {
		return errors.New("user_id is required")
	}
	if w.Currency == "" {
		return errors.New("currency is required")
	}
	if w.Balance < 0 {
		return errors.New("balance cannot be negative")
	}
	return nil
}

// Deposit adds amount to the wallet balance
func (w *Wallet) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be greater than 0")
	}
	w.Balance += amount
	return nil
}

// Withdraw subtracts amount from the wallet balance
func (w *Wallet) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be greater than 0")
	}
	if w.Balance < amount {
		return errors.New("insufficient balance")
	}
	w.Balance -= amount
	return nil
}

// WalletRepository defines the interface for wallet data access
type WalletRepository interface {
	Create(wallet *Wallet) error
	FindByID(id string) (*Wallet, error)
	FindByUserID(userID string) (*Wallet, error)
	Update(wallet *Wallet) error
}

// WalletService defines the interface for wallet business logic
type WalletService interface {
	CreateWallet(wallet *Wallet) error
	GetWalletByID(id string) (*Wallet, error)
	GetWalletByUserID(userID string) (*Wallet, error)
	UpdateWallet(wallet *Wallet) error
	GetWalletBalance(userID string) (float64, error)
} 