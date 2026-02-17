package pg

import (
	"context"
	"database/sql"
	"tracker-server/internal/repo"
)

type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) repo.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *repo.User, secret *repo.UserSecret) error {
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

func (r *postgresUserRepository) GetByID(ctx context.Context, id string) (*repo.User, error) {
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
