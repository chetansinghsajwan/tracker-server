package suite

import (
	"context"
	"testing"
	"time"

	"tracker-server/internal/repo"
)

type UserRepoSuite struct {
	Repo repo.UserRepo
}

func (s *UserRepoSuite) TestAll(t *testing.T) {
	t.Run("Create", s.TestCreate)
	t.Run("GetByID", s.TestGetByID)
}

func (s *UserRepoSuite) TestCreate(t *testing.T) {
	ctx := context.Background()
	user := &repo.User{
		ID:          "user-1",
		Email:       "test@example.com",
		FullName:    "Test User",
		DisplayName: "Test User",
		CreatedAt:   time.Now(),
	}
	secret := &repo.UserSecret{
		ID:    "user-1",
		Value: "hashed_secret",
	}

	err := s.Repo.Create(ctx, user, secret)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Verify creation implicitly via GetByID or direct DB check if possible?
	// Contract test usually relies on other methods of the interface.
	fetched, err := s.Repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID failed after Create: %v", err)
	}
	if fetched.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, fetched.Email)
	}
}

func (s *UserRepoSuite) TestGetByID(t *testing.T) {
	ctx := context.Background()
	// Assumes data exists or cleans up?
	// Tests should ideally be independent or run in a transaction that rolls back.
	// For simplicity in this suite, we assume the runner handles teardown or we create unique IDs.

	id := "user-2"
	user := &repo.User{
		ID:          id,
		Email:       "get@example.com",
		FullName:    "Get User",
		DisplayName: "Get User",
		CreatedAt:   time.Now(),
	}
	secret := &repo.UserSecret{
		ID:    id,
		Value: "secret",
	}

	if err := s.Repo.Create(ctx, user, secret); err != nil {
		t.Fatalf("Setup create failed: %v", err)
	}

	fetched, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if fetched.ID != id {
		t.Errorf("Expected ID %s, got %s", id, fetched.ID)
	}

	_, err = s.Repo.GetByID(ctx, "non-existent")
	if err == nil {
		t.Error("Expected error for non-existent ID, got nil")
	}
}
