package repo

import (
	"context"
	"time"
)



type Category struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // income, expense
	CreatedAt time.Time `json:"created_at"`
}



type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id int64) (*Category, error)
	ListByUserID(ctx context.Context, userID string) ([]*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id int64) error
}
