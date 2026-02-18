package suite

import (
	"context"
	"testing"
	"time"

	"tracker-server/internal/repo"
)

type TransactionRepoSuite struct {
	Repo         repo.TransactionRepo
	UserRepo     repo.UserRepo
	AccountRepo  repo.AccountRepo
	CategoryRepo repo.CategoryRepo
}

func (s *TransactionRepoSuite) TestAll(t *testing.T) {
	t.Run("CreateTransaction", s.TestCreate)
	t.Run("GetTransaction", s.TestGetByID)
	t.Run("ListTransactions", s.TestListByUserID)
}

func (s *TransactionRepoSuite) TestCreate(t *testing.T) {
	ctx := context.Background()
	userID := "tx-user-1"
	s.UserRepo.Create(ctx, &repo.User{ID: userID, Email: "tx1@ex.com", FullName: "Tx User", CreatedAt: time.Now()}, &repo.UserSecret{ID: userID, Value: "s"})

	acc1 := &repo.Account{UserID: userID, Name: "A1", Type: "cash", Currency: "USD", CreatedAt: time.Now()}
	s.AccountRepo.Create(ctx, acc1)
	cat1 := &repo.Category{UserID: userID, Name: "C1", Type: "expense", CreatedAt: time.Now()}
	s.CategoryRepo.Create(ctx, cat1)

	tx := &repo.Transaction{
		UserID:          userID,
		FromAccountID:   &acc1.ID,
		CategoryID:      &cat1.ID,
		Amount:          100.50,
		Type:            "expense",
		Description:     "Test Tx",
		TransactionDate: time.Now(),
		CreatedAt:       time.Now(),
	}

	if err := s.Repo.Create(ctx, tx); err != nil {
		t.Fatalf("Create tx failed: %v", err)
	}
	if tx.ID == 0 {
		t.Error("Expected ID to be set")
	}
}

func (s *TransactionRepoSuite) TestGetByID(t *testing.T) {
	ctx := context.Background()
	userID := "tx-user-2"
	s.UserRepo.Create(ctx, &repo.User{ID: userID, Email: "tx2@ex.com", FullName: "Tx User", CreatedAt: time.Now()}, &repo.UserSecret{ID: userID, Value: "s"})

	acc1 := &repo.Account{UserID: userID, Name: "A1", Type: "cash", Currency: "USD", CreatedAt: time.Now()}
	s.AccountRepo.Create(ctx, acc1)

	tx := &repo.Transaction{
		UserID:          userID,
		FromAccountID:   &acc1.ID,
		Amount:          50.00,
		Type:            "expense",
		Description:     "Get Tx",
		TransactionDate: time.Now(),
		CreatedAt:       time.Now(),
	}
	s.Repo.Create(ctx, tx)

	fetched, err := s.Repo.GetByID(ctx, tx.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if fetched.Amount != tx.Amount {
		t.Errorf("Expected amount %f, got %f", tx.Amount, fetched.Amount)
	}
}

func (s *TransactionRepoSuite) TestListByUserID(t *testing.T) {
	ctx := context.Background()
	userID := "tx-user-3"
	s.UserRepo.Create(ctx, &repo.User{ID: userID, Email: "tx3@ex.com", FullName: "Tx User", CreatedAt: time.Now()}, &repo.UserSecret{ID: userID, Value: "s"})
	acc1 := &repo.Account{UserID: userID, Name: "A1", Type: "cash", Currency: "USD", CreatedAt: time.Now()}
	s.AccountRepo.Create(ctx, acc1)

	s.Repo.Create(ctx, &repo.Transaction{UserID: userID, FromAccountID: &acc1.ID, Amount: 10, Type: "expense", TransactionDate: time.Now(), CreatedAt: time.Now()})
	s.Repo.Create(ctx, &repo.Transaction{UserID: userID, FromAccountID: &acc1.ID, Amount: 20, Type: "expense", TransactionDate: time.Now(), CreatedAt: time.Now()})

	txs, err := s.Repo.ListByUserID(ctx, userID, repo.TransactionFilter{})
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}
	if len(txs) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(txs))
	}
}
