package domain

import (
	"context"
	"errors"
	"time"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeTransfer   TransactionType = "transfer"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending     TransactionStatus = "pending"
	TransactionStatusCompleted   TransactionStatus = "completed"
	TransactionStatusFailed      TransactionStatus = "failed"
	TransactionStatusRolledBack  TransactionStatus = "rolled_back"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	FromUserID  *uint             `gorm:"index" json:"from_user_id,omitempty"`
	ToUserID    *uint             `gorm:"index" json:"to_user_id,omitempty"`
	Amount      float64           `gorm:"type:decimal(20,2);not null" json:"amount"`
	Type        TransactionType   `gorm:"size:20;not null" json:"type"`
	Status      TransactionStatus `gorm:"size:20;not null;default:pending" json:"status"`
	Description string            `gorm:"size:255" json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// Validate performs validation on the Transaction struct
func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if t.Type == "" {
		return errors.New("transaction type is required")
	}
	if t.Type == TransactionTypeTransfer && (t.FromUserID == nil || t.ToUserID == nil) {
		return errors.New("from_user_id and to_user_id are required for transfer transactions")
	}
	return nil
}

// CanTransitionTo checks if the transaction can transition to the new status
func (t *Transaction) CanTransitionTo(newStatus TransactionStatus) bool {
	switch t.Status {
	case TransactionStatusPending:
		return newStatus == TransactionStatusCompleted || newStatus == TransactionStatusFailed
	case TransactionStatusCompleted:
		return false
	case TransactionStatusFailed:
		return false
	default:
		return false
	}
}

// TransitionTo attempts to transition the transaction to a new status
func (t *Transaction) TransitionTo(newStatus TransactionStatus) error {
	if !t.CanTransitionTo(newStatus) {
		return errors.New("invalid status transition")
	}
	t.Status = newStatus
	return nil
}

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	FindByID(ctx context.Context, id uint) (*Transaction, error)
	FindByUserID(ctx context.Context, userID uint) ([]*Transaction, error)
	Update(ctx context.Context, transaction *Transaction) error
}

// TransactionService defines the interface for transaction business logic
type TransactionService interface {
	CreateTransaction(ctx context.Context, transaction *Transaction) error
	GetTransactionByID(ctx context.Context, id uint) (*Transaction, error)
	GetUserTransactions(ctx context.Context, userID uint) ([]*Transaction, error)
	ProcessTransaction(ctx context.Context, transaction *Transaction) error
} 