package repository

import "context"

func (r *UploadRepository) EnsureSchema(ctx context.Context) error {
	const q = `
CREATE TABLE IF NOT EXISTS uploads (
	id BIGSERIAL PRIMARY KEY,
	object_key TEXT NOT NULL UNIQUE,
	file_name TEXT NOT NULL,
	content_type TEXT NOT NULL,
	size BIGINT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`
	_, err := r.db.ExecContext(ctx, q)
	return err
}
