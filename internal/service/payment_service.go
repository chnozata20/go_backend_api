package service

import (
	"context"
	"errors"
	"time"

	"github.com/cihan-ozata/backend-path/internal/domain"
	"github.com/cihan-ozata/backend-path/internal/repository"
)

// PaymentService implements domain.PaymentService interface
type PaymentService struct {
	paymentRepo    *repository.PaymentRepository
	walletRepo     *repository.WalletRepository
	transactionRepo *repository.TransactionRepository
}

// NewPaymentService creates a new PaymentService instance
func NewPaymentService(
	paymentRepo *repository.PaymentRepository,
	walletRepo *repository.WalletRepository,
	transactionRepo *repository.TransactionRepository,
) *PaymentService {
	return &PaymentService{
		paymentRepo:    paymentRepo,
		walletRepo:     walletRepo,
		transactionRepo: transactionRepo,
	}
}

// CreatePayment creates a new payment
func (s *PaymentService) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	if err := payment.Validate(); err != nil {
		return err
	}

	// Set initial status
	payment.Status = domain.PaymentStatusPending

	return s.paymentRepo.Create(ctx, payment)
}

// GetPaymentByID retrieves a payment by ID
func (s *PaymentService) GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error) {
	return s.paymentRepo.FindByID(ctx, id)
}

// GetUserPayments retrieves all payments for a user
func (s *PaymentService) GetUserPayments(ctx context.Context, userID string) ([]*domain.Payment, error) {
	return s.paymentRepo.FindByUserID(ctx, userID)
}

// ProcessPayment processes a payment
func (s *PaymentService) ProcessPayment(ctx context.Context, payment *domain.Payment) error {
	if err := payment.Validate(); err != nil {
		return err
	}

	// Get user's wallet
	wallet, err := s.walletRepo.FindByUserID(ctx, payment.UserID)
	if err != nil {
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Process payment based on method
	switch payment.Method {
	case domain.PaymentMethodWallet:
		// Check if sufficient balance
		if wallet.Balance < payment.Amount {
			payment.Status = domain.PaymentStatusFailed
			return errors.New("insufficient balance")
		}

		// Create withdrawal transaction
		transaction := &domain.Transaction{
			FromUserID:  &payment.UserID,
			Type:        domain.TransactionTypeWithdrawal,
			Amount:      payment.Amount,
			Status:      domain.TransactionStatusCompleted,
			Description: payment.Description,
		}

		if err := s.transactionRepo.Create(ctx, transaction); err != nil {
			payment.Status = domain.PaymentStatusFailed
			return err
		}

		// Update wallet balance
		if err := wallet.Withdraw(payment.Amount); err != nil {
			payment.Status = domain.PaymentStatusFailed
			return err
		}

		if err := s.walletRepo.Update(ctx, wallet); err != nil {
			payment.Status = domain.PaymentStatusFailed
			return err
		}

	default:
		return errors.New("unsupported payment method")
	}

	// Update payment status
	payment.Status = domain.PaymentStatusCompleted
	payment.UpdatedAt = time.Now()

	return s.paymentRepo.Update(ctx, payment)
}

// RefundPayment refunds a payment
func (s *PaymentService) RefundPayment(ctx context.Context, paymentID string) error {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return err
	}
	if payment == nil {
		return errors.New("payment not found")
	}

	if payment.Status != domain.PaymentStatusCompleted {
		return errors.New("can only refund completed payments")
	}

	// Get user's wallet
	wallet, err := s.walletRepo.FindByUserID(ctx, payment.UserID)
	if err != nil {
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Create deposit transaction for refund
	transaction := &domain.Transaction{
		ToUserID:    &payment.UserID,
		Type:        domain.TransactionTypeDeposit,
		Amount:      payment.Amount,
		Status:      domain.TransactionStatusCompleted,
		Description: "Refund for payment: " + payment.ID,
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return err
	}

	// Update wallet balance
	if err := wallet.Deposit(payment.Amount); err != nil {
		return err
	}

	if err := s.walletRepo.Update(ctx, wallet); err != nil {
		return err
	}

	// Update payment status
	payment.Status = domain.PaymentStatusRefunded
	payment.UpdatedAt = time.Now()

	return s.paymentRepo.Update(ctx, payment)
}

// GetAllPayments returns all payments
func (s *PaymentService) GetAllPayments(ctx context.Context) ([]*domain.Payment, error) {
	return s.paymentRepo.FindAll(ctx)
} 