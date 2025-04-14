package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cihan-ozata/backend-path/internal/domain"
	"github.com/cihan-ozata/backend-path/internal/service"
)

type TransactionWorker struct {
	queue          chan *domain.Transaction
	workers        int
	wg             sync.WaitGroup
	stats          *TransactionStats
	transactionSvc *service.TransactionService
}

type TransactionStats struct {
	TotalProcessed    uint64
	TotalAmount       uint64 // Kuruş cinsinden tutulacak
	TotalErrors       uint64
	ProcessingTime    time.Duration
	LastProcessedTime time.Time
}

func NewTransactionWorker(workers int, transactionSvc *service.TransactionService) *TransactionWorker {
	return &TransactionWorker{
		queue:          make(chan *domain.Transaction, 1000),
		workers:        workers,
		stats:          &TransactionStats{},
		transactionSvc: transactionSvc,
	}
}

func (w *TransactionWorker) Start(ctx context.Context) {
	for i := 0; i < w.workers; i++ {
		w.wg.Add(1)
		go w.processTransactions(ctx)
	}
}

func (w *TransactionWorker) Stop() {
	close(w.queue)
	w.wg.Wait()
}

func (w *TransactionWorker) Enqueue(transaction *domain.Transaction) {
	w.queue <- transaction
}

func (w *TransactionWorker) processTransactions(ctx context.Context) {
	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case transaction, ok := <-w.queue:
			if !ok {
				return
			}
			startTime := time.Now()
			
			if err := w.transactionSvc.ProcessTransaction(ctx, transaction); err != nil {
				atomic.AddUint64(&w.stats.TotalErrors, 1)
				continue
			}

			atomic.AddUint64(&w.stats.TotalProcessed, 1)
			atomic.AddUint64(&w.stats.TotalAmount, uint64(transaction.Amount*100)) // Kuruş cinsinden ekle
			w.stats.ProcessingTime = time.Since(startTime)
			w.stats.LastProcessedTime = time.Now()
		}
	}
}

func (w *TransactionWorker) GetStats() TransactionStats {
	return *w.stats
} 