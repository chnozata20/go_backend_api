package repository

import (
	"context"
	"time"

	"github.com/cihan-ozata/backend-path/internal/domain"
	"gorm.io/gorm"
)

type BalanceRepository struct {
	db *gorm.DB
}

func NewBalanceRepository(db *gorm.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

func (r *BalanceRepository) Create(ctx context.Context, balance *domain.Balance) error {
	balance.CreatedAt = time.Now()
	balance.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(balance).Error
}

func (r *BalanceRepository) FindByUserID(ctx context.Context, userID uint) (*domain.Balance, error) {
	var balance domain.Balance
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&balance).Error
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (r *BalanceRepository) Update(ctx context.Context, balance *domain.Balance) error {
	balance.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(balance).Error
}

func (r *BalanceRepository) FindHistoryByUserID(ctx context.Context, userID uint, startDate, endDate time.Time) ([]*domain.Balance, error) {
	var balances []*domain.Balance
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND updated_at BETWEEN ? AND ?", userID, startDate, endDate).
		Order("updated_at DESC").
		Find(&balances).Error
	if err != nil {
		return nil, err
	}
	return balances, nil
}

func (r *BalanceRepository) FindTransactionsByUserID(ctx context.Context, userID uint) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	err := r.db.WithContext(ctx).
		Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Order("created_at DESC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
} 