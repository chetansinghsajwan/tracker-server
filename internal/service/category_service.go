package service

import (
	"context"
	"time"

	"tracker-server/internal/repo"
)

type CategoryService struct {
	categoryRepo repo.CategoryRepo
}

func NewCategoryService(categoryRepo repo.CategoryRepo) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) Create(ctx context.Context, userID, name, catType string) (*repo.Category, error) {
	category := &repo.Category{
		UserID:    userID,
		Name:      name,
		Type:      catType,
		CreatedAt: time.Now(),
	}
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id int64) (*repo.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

func (s *CategoryService) ListByUserID(ctx context.Context, userID string) ([]*repo.Category, error) {
	return s.categoryRepo.ListByUserID(ctx, userID)
}

func (s *CategoryService) Update(ctx context.Context, category *repo.Category) error {
	return s.categoryRepo.Update(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.categoryRepo.Delete(ctx, id)
}
