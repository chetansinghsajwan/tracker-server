package pg

import (
	"context"
	"database/sql"
	"tracker-server/internal/repo"
)

type postgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) repo.UserRepo {
	return &postgresUserRepo{db: db}
}

func (r *postgresUserRepo) Create(ctx context.Context, user *repo.User, secret *repo.UserSecret) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert User
	query := `INSERT INTO users (id, email, full_name, display_name, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, query, user.ID, user.Email, user.FullName, user.DisplayName, user.CreatedAt)
	if err != nil {
		return err
	}

	// Insert Secret
	querySecret := `INSERT INTO user_secrets (id, value) VALUES ($1, $2)`
	_, err = tx.ExecContext(ctx, querySecret, secret.ID, secret.Value)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *postgresUserRepo) GetByID(ctx context.Context, id string) (*repo.User, error) {
	query := `SELECT id, email, full_name, display_name, created_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var user repo.User
	var displayName sql.NullString // Handle nullable field

	err := row.Scan(&user.ID, &user.Email, &user.FullName, &displayName, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	if displayName.Valid {
		user.DisplayName = displayName.String
	}

	return &user, nil
}
