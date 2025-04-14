package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// PaymentMethod represents the payment method
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodWallet      PaymentMethod = "wallet"
)

// PaymentDetails represents payment-specific details
type PaymentDetails map[string]interface{}

// Value implements the driver.Valuer interface
func (pd PaymentDetails) Value() (driver.Value, error) {
	return json.Marshal(pd)
}

// Scan implements the sql.Scanner interface
func (pd *PaymentDetails) Scan(value interface{}) error {
	if value == nil {
		*pd = PaymentDetails{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, pd)
}

// Payment represents a payment transaction
type Payment struct {
	ID            string        `gorm:"primaryKey" json:"id"`
	UserID        uint          `gorm:"index" json:"user_id"`
	Amount        float64       `gorm:"type:decimal(20,2);not null" json:"amount"`
	Currency      string        `gorm:"size:3;not null" json:"currency"`
	Status        PaymentStatus `gorm:"size:20;not null" json:"status"`
	Method        PaymentMethod `gorm:"size:20;not null" json:"method"`
	Description   string        `gorm:"size:255" json:"description"`
	ReferenceID   string        `gorm:"size:100;uniqueIndex" json:"reference_id,omitempty"`
	PaymentDetails PaymentDetails `gorm:"type:jsonb" json:"payment_details,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name
func (Payment) TableName() string {
	return "payments"
}

// Validate performs validation on the Payment struct
func (p *Payment) Validate() error {
	if p.UserID == 0 {
		return errors.New("user_id is required")
	}
	if p.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if p.Currency == "" {
		return errors.New("currency is required")
	}
	if p.Status == "" {
		return errors.New("status is required")
	}
	if p.Method == "" {
		return errors.New("method is required")
	}
	return nil
}

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	Create(payment *Payment) error
	FindByID(id string) (*Payment, error)
	FindByUserID(userID string) ([]*Payment, error)
	Update(payment *Payment) error
}

// PaymentService defines the interface for payment business logic
type PaymentService interface {
	CreatePayment(payment *Payment) error
	GetPaymentByID(id string) (*Payment, error)
	GetUserPayments(userID string) ([]*Payment, error)
	ProcessPayment(payment *Payment) error
	RefundPayment(paymentID string) error
} 