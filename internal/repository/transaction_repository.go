package repository

import (
	"context"
	"errors"
	"time"

	"github.com/cihan-ozata/backend-path/internal/domain"
	"gorm.io/gorm"
)

// TransactionRepository implements domain.TransactionRepository interface
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new TransactionRepository instance
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(transaction).Error
}

// FindByID finds a transaction by ID
func (r *TransactionRepository) FindByID(ctx context.Context, id string) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.WithContext(ctx).First(&transaction, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

// FindByUserID finds transactions by user ID
func (r *TransactionRepository) FindByUserID(ctx context.Context, userID uint) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	err := r.db.WithContext(ctx).
		Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// Update updates a transaction
func (r *TransactionRepository) Update(ctx context.Context, transaction *domain.Transaction) error {
	transaction.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(transaction).Error
}

// FindAll returns all transactions
func (r *TransactionRepository) FindAll(ctx context.Context) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	if err := r.db.WithContext(ctx).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
} 