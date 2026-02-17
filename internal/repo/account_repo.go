package repo

import (
	"context"
	"time"
)



type Account struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}



type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	GetByID(ctx context.Context, id int64) (*Account, error)
	ListByUserID(ctx context.Context, userID string) ([]*Account, error)
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, id int64) error
}
