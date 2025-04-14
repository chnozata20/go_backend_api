package repository

import (
	"context"
	"errors"
	"time"

	"github.com/cihan-ozata/backend-path/internal/domain"
	"gorm.io/gorm"
)

// PaymentRepository implements domain.PaymentRepository interface
type PaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new PaymentRepository instance
func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create creates a new payment
func (r *PaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(payment).Error
}

// FindByID finds a payment by ID
func (r *PaymentRepository) FindByID(ctx context.Context, id string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).First(&payment, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// FindByUserID finds payments by user ID
func (r *PaymentRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

// Update updates a payment
func (r *PaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	payment.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(payment).Error
}

// FindAll returns all payments
func (r *PaymentRepository) FindAll(ctx context.Context) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	if err := r.db.WithContext(ctx).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
} 