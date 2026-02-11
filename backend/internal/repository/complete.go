package repository

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"upload/internal/entity"
)

type UploadRepository struct {
	db *sql.DB
}

func NewUploadRepository(db *sql.DB) *UploadRepository {
	return &UploadRepository{db: db}
}

func (r *UploadRepository) Create(ctx context.Context, u *entity.Upload) (int64, error) {
	const q = `
INSERT INTO uploads (object_key, file_name, content_type, size)
VALUES ($1, $2, $3, $4)
RETURNING id;
`
	var id int64
	err := r.db.QueryRowContext(ctx, q, u.ObjectKey, u.FileName, u.ContentType, u.Size).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
