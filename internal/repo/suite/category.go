package suite

import (
	"context"
	"testing"
	"time"

	"tracker-server/internal/repo"
)

type CategoryRepoSuite struct {
	Repo     repo.CategoryRepo
	UserRepo repo.UserRepo
}

func (s *CategoryRepoSuite) TestAll(t *testing.T) {
	t.Run("CreateCategory", s.TestCreate)
	t.Run("GetCategory", s.TestGetByID)
	t.Run("ListCategories", s.TestListByUserID)
}

func (s *CategoryRepoSuite) TestCreate(t *testing.T) {
	ctx := context.Background()
	userID := "cat-user-1"
	s.UserRepo.Create(ctx, &repo.User{ID: userID, Email: "cat1@ex.com", FullName: "Cat User", CreatedAt: time.Now()}, &repo.UserSecret{ID: userID, Value: "s"})

	category := &repo.Category{
		UserID:    userID,
		Name:      "Groceries",
		Type:      "expense",
		CreatedAt: time.Now(),
	}

	if err := s.Repo.Create(ctx, category); err != nil {
		t.Fatalf("Create category failed: %v", err)
	}
	if category.ID == 0 {
		t.Error("Expected ID to be set")
	}
}

func (s *CategoryRepoSuite) TestGetByID(t *testing.T) {
	ctx := context.Background()
	userID := "cat-user-2"
	s.UserRepo.Create(ctx, &repo.User{ID: userID, Email: "cat2@ex.com", FullName: "Cat User", CreatedAt: time.Now()}, &repo.UserSecret{ID: userID, Value: "s"})

	cat := &repo.Category{UserID: userID, Name: "Rent", Type: "expense", CreatedAt: time.Now()}
	s.Repo.Create(ctx, cat)

	fetched, err := s.Repo.GetByID(ctx, cat.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if fetched.Name != cat.Name {
		t.Errorf("Expected name %s, got %s", cat.Name, fetched.Name)
	}
}

func (s *CategoryRepoSuite) TestListByUserID(t *testing.T) {
	ctx := context.Background()
	userID := "cat-user-3"
	s.UserRepo.Create(ctx, &repo.User{ID: userID, Email: "cat3@ex.com", FullName: "Cat User", CreatedAt: time.Now()}, &repo.UserSecret{ID: userID, Value: "s"})

	s.Repo.Create(ctx, &repo.Category{UserID: userID, Name: "C1", Type: "expense", CreatedAt: time.Now()})
	s.Repo.Create(ctx, &repo.Category{UserID: userID, Name: "C2", Type: "income", CreatedAt: time.Now()})

	cats, err := s.Repo.ListByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}
	if len(cats) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(cats))
	}
}
