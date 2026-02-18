package suite

import (
	"context"
	"testing"
	"time"

	"tracker-server/internal/repo"
)

type AccountRepoSuite struct {
	Repo     repo.AccountRepo
	UserRepo repo.UserRepo // Needed to create dependency
}

func (s *AccountRepoSuite) TestAll(t *testing.T) {
	t.Run("CreateAccount", s.TestCreate)
	t.Run("GetAccount", s.TestGetByID)
	t.Run("ListAccounts", s.TestListByUserID)
}

func (s *AccountRepoSuite) TestCreate(t *testing.T) {
	ctx := context.Background()
	userID := "acc-user-1"

	// Setup user
	user := &repo.User{ID: userID, Email: "acc1@ex.com", FullName: "Acc User", DisplayName: "Acc User", CreatedAt: time.Now()}
	s.UserRepo.Create(ctx, user, &repo.UserSecret{ID: userID, Value: "s"})

	account := &repo.Account{
		UserID:    userID,
		Name:      "Test Account",
		Type:      "bank",
		Currency:  "USD",
		CreatedAt: time.Now(),
	}

	if err := s.Repo.Create(ctx, account); err != nil {
		t.Fatalf("Create account failed: %v", err)
	}
	if account.ID == 0 {
		t.Error("Expected ID to be set after Create")
	}
}

func (s *AccountRepoSuite) TestGetByID(t *testing.T) {
	ctx := context.Background()
	userID := "acc-user-2"
	user := &repo.User{ID: userID, Email: "acc2@ex.com", FullName: "Acc User", DisplayName: "Acc User", CreatedAt: time.Now()}
	s.UserRepo.Create(ctx, user, &repo.UserSecret{ID: userID, Value: "s"})

	account := &repo.Account{UserID: userID, Name: "Get Account", Type: "cash", Currency: "USD", CreatedAt: time.Now()}
	s.Repo.Create(ctx, account)

	fetched, err := s.Repo.GetByID(ctx, account.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if fetched.Name != account.Name {
		t.Errorf("Expected name %s, got %s", account.Name, fetched.Name)
	}
}

func (s *AccountRepoSuite) TestListByUserID(t *testing.T) {
	ctx := context.Background()
	userID := "acc-user-3"
	user := &repo.User{ID: userID, Email: "acc3@ex.com", FullName: "Acc User", DisplayName: "Acc User", CreatedAt: time.Now()}
	s.UserRepo.Create(ctx, user, &repo.UserSecret{ID: userID, Value: "s"})

	s.Repo.Create(ctx, &repo.Account{UserID: userID, Name: "A1", Type: "bank", Currency: "USD", CreatedAt: time.Now()})
	s.Repo.Create(ctx, &repo.Account{UserID: userID, Name: "A2", Type: "cash", Currency: "EUR", CreatedAt: time.Now()})

	accounts, err := s.Repo.ListByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}
	if len(accounts) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(accounts))
	}
}
