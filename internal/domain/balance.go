package domain

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Balance represents a user's balance
type Balance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;unique" json:"user_id"`
	Amount    float64   `gorm:"type:decimal(20,2);not null;default:0" json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	mu        sync.RWMutex
}

// Validate performs validation on the Balance struct
func (b *Balance) Validate() error {
	if b.UserID == 0 {
		return errors.New("user_id is required")
	}
	return nil
}

// Add adds amount to the balance
func (b *Balance) Add(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.Amount += amount
	return nil
}

// Subtract subtracts amount from the balance
func (b *Balance) Subtract(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.Amount < amount {
		return errors.New("insufficient balance")
	}

	b.Amount -= amount
	return nil
}

// GetAmount returns the current balance amount
func (b *Balance) GetAmount() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Amount
}

// BalanceRepository defines the interface for balance data access
type BalanceRepository interface {
	Create(ctx context.Context, balance *Balance) error
	FindByUserID(ctx context.Context, userID uint) (*Balance, error)
	Update(ctx context.Context, balance *Balance) error
}

// BalanceService defines the interface for balance business logic
type BalanceService interface {
	GetBalance(ctx context.Context, userID uint) (*Balance, error)
	AddBalance(ctx context.Context, userID uint, amount float64) error
	SubtractBalance(ctx context.Context, userID uint, amount float64) error
	TransferBalance(ctx context.Context, fromUserID, toUserID uint, amount float64) error
} 