package service

import (
	"context"
	"errors"
	"time"

	"github.com/cihan-ozata/backend-path/internal/domain"
	"github.com/cihan-ozata/backend-path/internal/repository"
)

// TransactionService implements domain.TransactionService interface
type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	walletRepo     *repository.WalletRepository
}

// NewTransactionService creates a new TransactionService instance
func NewTransactionService(transactionRepo *repository.TransactionRepository, walletRepo *repository.WalletRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		walletRepo:     walletRepo,
	}
}

// CreateTransaction creates a new transaction
func (s *TransactionService) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	if err := transaction.Validate(); err != nil {
		return err
	}

	// Get user's wallet
	wallet, err := s.walletRepo.FindByUserID(ctx, *transaction.FromUserID)
	if err != nil {
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Process transaction based on type
	switch transaction.Type {
	case domain.TransactionTypeDeposit:
		if err := wallet.Deposit(transaction.Amount); err != nil {
			return err
		}
	case domain.TransactionTypeWithdrawal:
		if err := wallet.Withdraw(transaction.Amount); err != nil {
			return err
		}
	case domain.TransactionTypeTransfer:
		// Get recipient's wallet
		recipientWallet, err := s.walletRepo.FindByUserID(ctx, *transaction.ToUserID)
		if err != nil {
			return err
		}
		if recipientWallet == nil {
			return errors.New("recipient wallet not found")
		}

		// Withdraw from sender
		if err := wallet.Withdraw(transaction.Amount); err != nil {
			return err
		}

		// Deposit to recipient
		if err := recipientWallet.Deposit(transaction.Amount); err != nil {
			// Rollback sender's withdrawal
			if err := wallet.Deposit(transaction.Amount); err != nil {
				return err
			}
			return err
		}

		// Update recipient's wallet
		if err := s.walletRepo.Update(ctx, recipientWallet); err != nil {
			// Rollback both operations
			if err := wallet.Deposit(transaction.Amount); err != nil {
				return err
			}
			if err := recipientWallet.Withdraw(transaction.Amount); err != nil {
				return err
			}
			return err
		}
	default:
		return errors.New("invalid transaction type")
	}

	// Update wallet balance
	if err := s.walletRepo.Update(ctx, wallet); err != nil {
		return err
	}

	// Set transaction status
	transaction.Status = domain.TransactionStatusCompleted
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	return s.transactionRepo.Create(ctx, transaction)
}

// RollbackTransaction rolls back a completed transaction
func (s *TransactionService) RollbackTransaction(ctx context.Context, transactionID string) error {
	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return err
	}
	if transaction == nil {
		return errors.New("transaction not found")
	}

	if transaction.Status != domain.TransactionStatusCompleted {
		return errors.New("can only rollback completed transactions")
	}

	// Get user's wallet
	wallet, err := s.walletRepo.FindByUserID(ctx, *transaction.FromUserID)
	if err != nil {
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Reverse transaction based on type
	switch transaction.Type {
	case domain.TransactionTypeDeposit:
		if err := wallet.Withdraw(transaction.Amount); err != nil {
			return err
		}
	case domain.TransactionTypeWithdrawal:
		if err := wallet.Deposit(transaction.Amount); err != nil {
			return err
		}
	case domain.TransactionTypeTransfer:
		// Get recipient's wallet
		recipientWallet, err := s.walletRepo.FindByUserID(ctx, *transaction.ToUserID)
		if err != nil {
			return err
		}
		if recipientWallet == nil {
			return errors.New("recipient wallet not found")
		}

		// Deposit back to sender
		if err := wallet.Deposit(transaction.Amount); err != nil {
			return err
		}

		// Withdraw from recipient
		if err := recipientWallet.Withdraw(transaction.Amount); err != nil {
			// Rollback sender's deposit
			if err := wallet.Withdraw(transaction.Amount); err != nil {
				return err
			}
			return err
		}

		// Update recipient's wallet
		if err := s.walletRepo.Update(ctx, recipientWallet); err != nil {
			// Rollback both operations
			if err := wallet.Withdraw(transaction.Amount); err != nil {
				return err
			}
			if err := recipientWallet.Deposit(transaction.Amount); err != nil {
				return err
			}
			return err
		}
	default:
		return errors.New("invalid transaction type")
	}

	// Update wallet balance
	if err := s.walletRepo.Update(ctx, wallet); err != nil {
		return err
	}

	// Update transaction status
	transaction.Status = domain.TransactionStatusRolledBack
	transaction.UpdatedAt = time.Now()

	return s.transactionRepo.Update(ctx, transaction)
}

// GetTransactionByID retrieves a transaction by ID
func (s *TransactionService) GetTransactionByID(ctx context.Context, id string) (*domain.Transaction, error) {
	return s.transactionRepo.FindByID(ctx, id)
}

// GetUserTransactions retrieves all transactions for a user
func (s *TransactionService) GetUserTransactions(ctx context.Context, userID uint) ([]*domain.Transaction, error) {
	return s.transactionRepo.FindByUserID(ctx, userID)
}

// ProcessTransaction processes a transaction
func (s *TransactionService) ProcessTransaction(ctx context.Context, transaction *domain.Transaction) error {
	if err := transaction.Validate(); err != nil {
		return err
	}

	// Get user's wallet
	wallet, err := s.walletRepo.FindByUserID(ctx, *transaction.FromUserID)
	if err != nil {
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Process transaction based on type
	switch transaction.Type {
	case domain.TransactionTypeDeposit:
		if err := wallet.Deposit(transaction.Amount); err != nil {
			return err
		}
	case domain.TransactionTypeWithdrawal:
		if err := wallet.Withdraw(transaction.Amount); err != nil {
			return err
		}
	default:
		return errors.New("invalid transaction type")
	}

	// Update wallet balance
	if err := s.walletRepo.Update(ctx, wallet); err != nil {
		return err
	}

	// Update transaction status
	transaction.Status = domain.TransactionStatusCompleted
	transaction.UpdatedAt = time.Now()

	return s.transactionRepo.Update(ctx, transaction)
}

// GetAllTransactions returns all transactions
func (s *TransactionService) GetAllTransactions(ctx context.Context) ([]*domain.Transaction, error) {
	return s.transactionRepo.FindAll(ctx)
} 