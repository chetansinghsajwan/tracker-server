package service

import (
	"context"
	"time"

	"tracker-server/internal/repo"
)

type TagService struct {
	tagRepo repo.TagRepo
}

func NewTagService(tagRepo repo.TagRepo) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

func (s *TagService) Create(ctx context.Context, userID, name string) (*repo.Tag, error) {
	tag := &repo.Tag{
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
	}
	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagService) GetByID(ctx context.Context, id int64) (*repo.Tag, error) {
	return s.tagRepo.GetByID(ctx, id)
}

func (s *TagService) ListByUserID(ctx context.Context, userID string) ([]*repo.Tag, error) {
	return s.tagRepo.ListByUserID(ctx, userID)
}

func (s *TagService) Update(ctx context.Context, tag *repo.Tag) error {
	return s.tagRepo.Update(ctx, tag)
}

func (s *TagService) Delete(ctx context.Context, id int64) error {
	return s.tagRepo.Delete(ctx, id)
}
