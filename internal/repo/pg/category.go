package pg

import (
	"context"
	"database/sql"
	"tracker-server/internal/repo"
)

type postgresCategoryRepository struct {
	db *sql.DB
}

func NewPostgresCategoryRepository(db *sql.DB) repo.CategoryRepository {
	return &postgresCategoryRepository{db: db}
}

func (r *postgresCategoryRepository) Create(ctx context.Context, category *repo.Category) error {
	query := `INSERT INTO categories (user_id, name, type, created_at)
			  VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowContext(ctx, query, category.UserID, category.Name, category.Type, category.CreatedAt).Scan(&category.ID)
}

func (r *postgresCategoryRepository) GetByID(ctx context.Context, id int64) (*repo.Category, error) {
	query := `SELECT id, user_id, name, type, created_at FROM categories WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var category repo.Category
	err := row.Scan(&category.ID, &category.UserID, &category.Name, &category.Type, &category.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *postgresCategoryRepository) ListByUserID(ctx context.Context, userID string) ([]*repo.Category, error) {
	query := `SELECT id, user_id, name, type, created_at FROM categories WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*repo.Category
	for rows.Next() {
		var category repo.Category
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Type, &category.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	return categories, nil
}

func (r *postgresCategoryRepository) Update(ctx context.Context, category *repo.Category) error {
	query := `UPDATE categories SET name = $1, type = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, category.Name, category.Type, category.ID)
	return err
}

func (r *postgresCategoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
