package service

import (
	"context"
	"time"

	"tracker-server/internal/repo"
)

type AccountService struct {
	accountRepo repo.AccountRepo
}

func NewAccountService(accountRepo repo.AccountRepo) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

func (s *AccountService) Create(ctx context.Context, userID, name, acctType, currency string) (*repo.Account, error) {
	account := &repo.Account{
		UserID:    userID,
		Name:      name,
		Type:      acctType,
		Currency:  currency,
		CreatedAt: time.Now(),
	}
	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) GetByID(ctx context.Context, id int64) (*repo.Account, error) {
	return s.accountRepo.GetByID(ctx, id)
}

func (s *AccountService) ListByUserID(ctx context.Context, userID string) ([]*repo.Account, error) {
	return s.accountRepo.ListByUserID(ctx, userID)
}

func (s *AccountService) Update(ctx context.Context, account *repo.Account) error {
	return s.accountRepo.Update(ctx, account)
}

func (s *AccountService) Delete(ctx context.Context, id int64) error {
	return s.accountRepo.Delete(ctx, id)
}
