package pg

import (
	"context"
	"database/sql"
	"tracker-server/internal/repo"
)

type postgresAccountRepo struct {
	db *sql.DB
}

func NewPostgresAccountRepo(db *sql.DB) repo.AccountRepo {
	return &postgresAccountRepo{db: db}
}

func (r *postgresAccountRepo) Create(ctx context.Context, account *repo.Account) error {
	query := `INSERT INTO accounts (user_id, name, type, currency, created_at)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRowContext(ctx, query, account.UserID, account.Name, account.Type, account.Currency, account.CreatedAt).Scan(&account.ID)
}

func (r *postgresAccountRepo) GetByID(ctx context.Context, id int64) (*repo.Account, error) {
	query := `SELECT id, user_id, name, type, currency, created_at FROM accounts WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var account repo.Account
	err := row.Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.Currency, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *postgresAccountRepo) ListByUserID(ctx context.Context, userID string) ([]*repo.Account, error) {
	query := `SELECT id, user_id, name, type, currency, created_at FROM accounts WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*repo.Account
	for rows.Next() {
		var account repo.Account
		if err := rows.Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.Currency, &account.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

func (r *postgresAccountRepo) Update(ctx context.Context, account *repo.Account) error {
	query := `UPDATE accounts SET name = $1, type = $2, currency = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, account.Name, account.Type, account.Currency, account.ID)
	return err
}

func (r *postgresAccountRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
