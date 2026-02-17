package pg

import (
	"context"
	"database/sql"
	"tracker-server/internal/repo"
)

type postgresTagRepository struct {
	db *sql.DB
}

func NewPostgresTagRepository(db *sql.DB) repo.TagRepository {
	return &postgresTagRepository{db: db}
}

func (r *postgresTagRepository) Create(ctx context.Context, tag *repo.Tag) error {
	query := `INSERT INTO tags (user_id, name, created_at)
			  VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, tag.UserID, tag.Name, tag.CreatedAt).Scan(&tag.ID)
}

func (r *postgresTagRepository) GetByID(ctx context.Context, id int64) (*repo.Tag, error) {
	query := `SELECT id, user_id, name, created_at FROM tags WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var tag repo.Tag
	err := row.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *postgresTagRepository) ListByUserID(ctx context.Context, userID string) ([]*repo.Tag, error) {
	query := `SELECT id, user_id, name, created_at FROM tags WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*repo.Tag
	for rows.Next() {
		var tag repo.Tag
		if err := rows.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}

func (r *postgresTagRepository) Update(ctx context.Context, tag *repo.Tag) error {
	query := `UPDATE tags SET name = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, tag.Name, tag.ID)
	return err
}

func (r *postgresTagRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tags WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
