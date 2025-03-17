package models

import (
	"time"

	"gorm.io/gorm"
)

// TransactionType, işlem tiplerini temsil eder
type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeTransfer TransactionType = "transfer"
)

// TransactionStatus, işlem durumlarını temsil eder
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

// Transaction, işlem modelini temsil eder
type Transaction struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	FromUserID *uint             `gorm:"index" json:"from_user_id,omitempty"` // Deposit işlemlerinde NULL olabilir
	ToUserID   *uint             `gorm:"index" json:"to_user_id,omitempty"`   // Withdraw işlemlerinde NULL olabilir
	Amount     float64           `gorm:"type:decimal(20,2);not null" json:"amount"`
	Type       TransactionType   `gorm:"size:20;not null" json:"type"`
	Status     TransactionStatus `gorm:"size:20;not null;default:pending" json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	DeletedAt  gorm.DeletedAt    `gorm:"index" json:"-"` // Soft delete için

	// İlişkiler
	FromUser *User `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUser   *User `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
}

// TableName, tablonun adını belirtir
func (Transaction) TableName() string {
	return "transactions"
} 