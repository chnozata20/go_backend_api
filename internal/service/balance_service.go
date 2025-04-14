package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/cihan-ozata/backend-path/internal/domain"
	"github.com/cihan-ozata/backend-path/internal/repository"
)

type BalanceService struct {
	balanceRepo *repository.BalanceRepository
	mu          sync.RWMutex
	cache       map[uint]*domain.Balance
}

func NewBalanceService(balanceRepo *repository.BalanceRepository) *BalanceService {
	return &BalanceService{
		balanceRepo: balanceRepo,
		cache:       make(map[uint]*domain.Balance),
	}
}

// GetBalance retrieves the current balance for a user
func (s *BalanceService) GetBalance(ctx context.Context, userID uint) (*domain.Balance, error) {
	s.mu.RLock()
	if balance, ok := s.cache[userID]; ok {
		s.mu.RUnlock()
		return balance, nil
	}
	s.mu.RUnlock()

	balance, err := s.balanceRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if balance == nil {
		return nil, errors.New("balance not found")
	}

	s.mu.Lock()
	s.cache[userID] = balance
	s.mu.Unlock()

	return balance, nil
}

// UpdateBalance updates the balance for a user
func (s *BalanceService) UpdateBalance(ctx context.Context, balance *domain.Balance) error {
	if err := balance.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.balanceRepo.Update(ctx, balance); err != nil {
		return err
	}

	s.cache[balance.UserID] = balance
	return nil
}

// GetBalanceHistory retrieves balance history for a user
func (s *BalanceService) GetBalanceHistory(ctx context.Context, userID uint, startDate, endDate time.Time) ([]*domain.Balance, error) {
	return s.balanceRepo.FindHistoryByUserID(ctx, userID, startDate, endDate)
}

// CalculateBalance calculates the current balance based on transaction history
func (s *BalanceService) CalculateBalance(ctx context.Context, userID uint) (*domain.Balance, error) {
	// Get cached balance
	s.mu.RLock()
	if balance, ok := s.cache[userID]; ok {
		s.mu.RUnlock()
		return balance, nil
	}
	s.mu.RUnlock()

	// Calculate from transactions if not in cache
	transactions, err := s.balanceRepo.FindTransactionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var totalAmount float64
	for _, tx := range transactions {
		switch tx.Type {
		case domain.TransactionTypeDeposit:
			totalAmount += tx.Amount
		case domain.TransactionTypeWithdrawal:
			totalAmount -= tx.Amount
		}
	}

	balance := &domain.Balance{
		UserID:    userID,
		Amount:    totalAmount,
		UpdatedAt: time.Now(),
	}

	s.mu.Lock()
	s.cache[userID] = balance
	s.mu.Unlock()

	return balance, nil
}

// ClearCache clears the balance cache
func (s *BalanceService) ClearCache() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = make(map[uint]*domain.Balance)
} 