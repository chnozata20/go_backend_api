package service

import (
	"context"
	"errors"

	"github.com/cihan-ozata/backend-path/internal/domain"
	"github.com/cihan-ozata/backend-path/internal/repository"
)

// WalletService implements domain.WalletService interface
type WalletService struct {
	walletRepo *repository.WalletRepository
}

// NewWalletService creates a new WalletService instance
func NewWalletService(walletRepo *repository.WalletRepository) *WalletService {
	return &WalletService{walletRepo: walletRepo}
}

// CreateWallet creates a new wallet
func (s *WalletService) CreateWallet(ctx context.Context, wallet *domain.Wallet) error {
	if err := wallet.Validate(); err != nil {
		return err
	}

	// Check if wallet already exists for user
	existingWallet, err := s.walletRepo.FindByUserID(ctx, wallet.UserID)
	if err != nil {
		return err
	}
	if existingWallet != nil {
		return errors.New("wallet already exists for user")
	}

	return s.walletRepo.Create(ctx, wallet)
}

// GetWalletByID retrieves a wallet by ID
func (s *WalletService) GetWalletByID(ctx context.Context, id string) (*domain.Wallet, error) {
	return s.walletRepo.FindByID(ctx, id)
}

// GetWallet retrieves a wallet by user ID
func (s *WalletService) GetWallet(ctx context.Context, userID uint) (*domain.Wallet, error) {
	return s.walletRepo.FindByUserID(ctx, userID)
}

// UpdateWallet updates a wallet
func (s *WalletService) UpdateWallet(ctx context.Context, wallet *domain.Wallet) error {
	if err := wallet.Validate(); err != nil {
		return err
	}

	return s.walletRepo.Update(ctx, wallet)
}

// GetWalletBalance retrieves the balance for a user's wallet
func (s *WalletService) GetWalletBalance(ctx context.Context, userID uint) (float64, error) {
	wallet, err := s.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	if wallet == nil {
		return 0, errors.New("wallet not found")
	}
	return wallet.Balance, nil
}

// DeleteWallet deletes a wallet by user ID
func (s *WalletService) DeleteWallet(ctx context.Context, userID uint) error {
	return s.walletRepo.Delete(ctx, userID)
}

// GetWalletByUserID retrieves a wallet by user ID
func (s *WalletService) GetWalletByUserID(ctx context.Context, userID uint) (*domain.Wallet, error) {
	return s.walletRepo.FindByUserID(ctx, userID)
}

// GetAllWallets retrieves all wallets
func (s *WalletService) GetAllWallets(ctx context.Context) ([]*domain.Wallet, error) {
	return s.walletRepo.FindAll(ctx)
} 