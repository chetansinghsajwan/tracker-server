package suite

import (
	"context"
	"testing"
	"time"

	"tracker-server/internal/repo"
)

type TagRepoSuite struct {
	Repo     repo.TagRepo
	UserRepo repo.UserRepo
}

func (s *TagRepoSuite) TestAll(t *testing.T) {
	t.Run("CreateTag", s.TestCreate)
	t.Run("GetTag", s.TestGetByID)
	t.Run("ListTags", s.TestListByUserID)
}

func (s *TagRepoSuite) TestCreate(t *testing.T) {
	ctx := context.Background()
	userID := "tag-user-1"
	s.UserRepo.Create(ctx, &repo.User{ID: userID, Email: "tag1@ex.com", FullName: "Tag User", CreatedAt: time.Now()}, &repo.UserSecret{ID: userID, Value: "s"})

	tag := &repo.Tag{
		UserID:    userID,
		Name:      "urgent",
		CreatedAt: time.Now(),
	}

	if err := s.Repo.Create(ctx, tag); err != nil {
		t.Fatalf("Create tag failed: %v", err)
	}
	if tag.ID == 0 {
		t.Error("Expected ID to be set")
	}
}

func (s *TagRepoSuite) TestGetByID(t *testing.T) {
	ctx := context.Background()
	userID := "tag-user-2"
	s.UserRepo.Create(ctx, &repo.User{ID: userID, Email: "tag2@ex.com", FullName: "Tag User", CreatedAt: time.Now()}, &repo.UserSecret{ID: userID, Value: "s"})

	tag := &repo.Tag{UserID: userID, Name: "later", CreatedAt: time.Now()}
	s.Repo.Create(ctx, tag)

	fetched, err := s.Repo.GetByID(ctx, tag.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if fetched.Name != tag.Name {
		t.Errorf("Expected name %s, got %s", tag.Name, fetched.Name)
	}
}

func (s *TagRepoSuite) TestListByUserID(t *testing.T) {
	ctx := context.Background()
	userID := "tag-user-3"
	s.UserRepo.Create(ctx, &repo.User{ID: userID, Email: "tag3@ex.com", FullName: "Tag User", CreatedAt: time.Now()}, &repo.UserSecret{ID: userID, Value: "s"})

	s.Repo.Create(ctx, &repo.Tag{UserID: userID, Name: "T1", CreatedAt: time.Now()})
	s.Repo.Create(ctx, &repo.Tag{UserID: userID, Name: "T2", CreatedAt: time.Now()})

	tags, err := s.Repo.ListByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}
}
