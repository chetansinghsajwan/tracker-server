package service

import (
	"context"
	"time"

	"tracker-server/internal/repo"
)

type TransactionService struct {
	transactionRepo repo.TransactionRepo
}

func NewTransactionService(transactionRepo repo.TransactionRepo) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *TransactionService) Create(ctx context.Context, t *repo.Transaction) error {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	return s.transactionRepo.Create(ctx, t)
}

func (s *TransactionService) GetByID(ctx context.Context, id int64) (*repo.Transaction, error) {
	return s.transactionRepo.GetByID(ctx, id)
}

func (s *TransactionService) ListByUserID(ctx context.Context, userID string, filter repo.TransactionFilter) ([]*repo.Transaction, error) {
	return s.transactionRepo.ListByUserID(ctx, userID, filter)
}

func (s *TransactionService) Update(ctx context.Context, t *repo.Transaction) error {
	return s.transactionRepo.Update(ctx, t)
}

func (s *TransactionService) Delete(ctx context.Context, id int64) error {
	return s.transactionRepo.Delete(ctx, id)
}

func (s *TransactionService) AddTag(ctx context.Context, transactionID, tagID int64) error {
	return s.transactionRepo.AddTag(ctx, transactionID, tagID)
}

func (s *TransactionService) RemoveTag(ctx context.Context, transactionID, tagID int64) error {
	return s.transactionRepo.RemoveTag(ctx, transactionID, tagID)
}
