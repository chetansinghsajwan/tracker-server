package pg

import (
	"context"
	"database/sql"
	"fmt"
	"tracker-server/internal/repo"

	"github.com/lib/pq"
)

type postgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) repo.TransactionRepository {
	return &postgresTransactionRepository{db: db}
}

func (r *postgresTransactionRepository) Create(ctx context.Context, t *repo.Transaction) error {
	query := `INSERT INTO transactions
		(user_id, from_account_id, to_account_id, category_id, amount, type, description, transaction_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`
	return r.db.QueryRowContext(ctx, query,
		t.UserID, t.FromAccountID, t.ToAccountID, t.CategoryID, t.Amount, t.Type, t.Description, t.TransactionDate, t.CreatedAt,
	).Scan(&t.ID)
}

func (r *postgresTransactionRepository) GetByID(ctx context.Context, id int64) (*repo.Transaction, error) {
	query := `SELECT id, user_id, from_account_id, to_account_id, category_id, amount, type, description, transaction_date, created_at
			  FROM transactions WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var t repo.Transaction
	var desc sql.NullString // Description is nullable
	err := row.Scan(
		&t.ID, &t.UserID, &t.FromAccountID, &t.ToAccountID, &t.CategoryID, &t.Amount, &t.Type, &desc, &t.TransactionDate, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		t.Description = desc.String
	}
	return &t, nil
}

func (r *postgresTransactionRepository) ListByUserID(ctx context.Context, userID string, filter repo.TransactionFilter) ([]*repo.Transaction, error) {
	query := `SELECT t.id, t.user_id, t.from_account_id, t.to_account_id, t.category_id, t.amount, t.type, t.description, t.transaction_date, t.created_at
			  FROM transactions t WHERE t.user_id = $1`
	args := []interface{}{userID}
	argIdx := 2

	if filter.FromAccountID != nil {
		query += fmt.Sprintf(" AND t.from_account_id = $%d", argIdx)
		args = append(args, *filter.FromAccountID)
		argIdx++
	}
	if filter.ToAccountID != nil {
		query += fmt.Sprintf(" AND t.to_account_id = $%d", argIdx)
		args = append(args, *filter.ToAccountID)
		argIdx++
	}
	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND t.category_id = $%d", argIdx)
		args = append(args, *filter.CategoryID)
		argIdx++
	}
	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND t.transaction_date >= $%d", argIdx)
		args = append(args, *filter.StartDate)
		argIdx++
	}
	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND t.transaction_date <= $%d", argIdx)
		args = append(args, *filter.EndDate)
		argIdx++
	}
	if filter.MinAmount != nil {
		query += fmt.Sprintf(" AND t.amount >= $%d", argIdx)
		args = append(args, *filter.MinAmount)
		argIdx++
	}
	if filter.MaxAmount != nil {
		query += fmt.Sprintf(" AND t.amount <= $%d", argIdx)
		args = append(args, *filter.MaxAmount)
		argIdx++
	}
	// Tag filtering (assuming filter.Tags contains tag names)
	if len(filter.Tags) > 0 {
		// This requires a JOIN or EXISTS subquery
		// Using EXISTS for cleaner filtering logic
		// AND EXISTS (SELECT 1 FROM transaction_tags tt JOIN tags tg ON tt.tag_id = tg.id WHERE tt.transaction_id = t.id AND tg.name IN (...))
		query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM transaction_tags tt JOIN tags tg ON tt.tag_id = tg.id WHERE tt.transaction_id = t.id AND tg.name = ANY($%d))", argIdx)
		// PostgreSQL ANY takes an array
		args = append(args, pq.Array(filter.Tags))
		argIdx++
	}

	query += ` ORDER BY t.transaction_date DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*repo.Transaction
	for rows.Next() {
		var t repo.Transaction
		var desc sql.NullString
		if err := rows.Scan(&t.ID, &t.UserID, &t.FromAccountID, &t.ToAccountID, &t.CategoryID, &t.Amount, &t.Type, &desc, &t.TransactionDate, &t.CreatedAt); err != nil {
			return nil, err
		}
		if desc.Valid {
			t.Description = desc.String
		}
		transactions = append(transactions, &t)
	}
	return transactions, nil
}

func (r *postgresTransactionRepository) Update(ctx context.Context, t *repo.Transaction) error {
	query := `UPDATE transactions
		SET from_account_id=$1, to_account_id=$2, category_id=$3, amount=$4, type=$5, description=$6, transaction_date=$7
		WHERE id=$8`
	_, err := r.db.ExecContext(ctx, query,
		t.FromAccountID, t.ToAccountID, t.CategoryID, t.Amount, t.Type, t.Description, t.TransactionDate, t.ID,
	)
	return err
}

func (r *postgresTransactionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM transactions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresTransactionRepository) AddTag(ctx context.Context, transactionID int64, tagID int64) error {
	query := `INSERT INTO transaction_tags (transaction_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, transactionID, tagID)
	return err
}

func (r *postgresTransactionRepository) RemoveTag(ctx context.Context, transactionID int64, tagID int64) error {
	query := `DELETE FROM transaction_tags WHERE transaction_id = $1 AND tag_id = $2`
	_, err := r.db.ExecContext(ctx, query, transactionID, tagID)
	return err
}
