package repo

import (
	"context"
	"time"
)



type Transaction struct {
	ID              int64     `json:"id"`
	UserID          string    `json:"user_id"`
	FromAccountID   *int64    `json:"from_account_id,omitempty"`
	ToAccountID     *int64    `json:"to_account_id,omitempty"`
	CategoryID      *int64    `json:"category_id,omitempty"`
	Amount          float64   `json:"amount"`
	Type            string    `json:"type"` // income, expense, transfer
	Description     string    `json:"description,omitempty"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
	Tags            []Tag     `json:"tags,omitempty"` // For preloading tags
}

type TransactionFilter struct {
	FromAccountID *int64
	ToAccountID   *int64
	CategoryID    *int64
	StartDate     *time.Time
	EndDate       *time.Time
	MinAmount     *float64
	MaxAmount     *float64
	Tags          []string
}



type TransactionRepo interface {
	Create(ctx context.Context, transaction *Transaction) error
	GetByID(ctx context.Context, id int64) (*Transaction, error)
	ListByUserID(ctx context.Context, userID string, filter TransactionFilter) ([]*Transaction, error)
	Update(ctx context.Context, transaction *Transaction) error
	Delete(ctx context.Context, id int64) error
	AddTag(ctx context.Context, transactionID int64, tagID int64) error
	RemoveTag(ctx context.Context, transactionID int64, tagID int64) error
}
