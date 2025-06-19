package repository

import (
	"context"
	"time"

	"github.com/cihan-ozata/backend-path/internal/domain"
	"gorm.io/gorm"
)

// WalletRepository implements domain.WalletRepository interface
type WalletRepository struct {
	db *gorm.DB
}

// NewWalletRepository creates a new WalletRepository instance
func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// Create creates a new wallet
func (r *WalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(wallet).Error
}

// FindByID finds a wallet by ID
func (r *WalletRepository) FindByID(ctx context.Context, id string) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.WithContext(ctx).First(&wallet, id).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// FindByUserID finds a wallet by user ID
func (r *WalletRepository) FindByUserID(ctx context.Context, userID uint) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// Update updates a wallet
func (r *WalletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {
	wallet.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(wallet).Error
}

// Delete deletes a wallet by user ID
func (r *WalletRepository) Delete(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.Wallet{}).Error
}

// FindAll finds all wallets
func (r *WalletRepository) FindAll(ctx context.Context) ([]*domain.Wallet, error) {
	var wallets []*domain.Wallet
	err := r.db.WithContext(ctx).Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, nil
} 