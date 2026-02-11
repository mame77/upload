package port

import (
	"context"
	"time"

	"upload/internal/entity"
)

type UploadRepository interface {
	EnsureSchema(ctx context.Context) error
	Create(ctx context.Context, u *entity.Upload) (int64, error)
}

type PresignService interface {
	PresignPut(ctx context.Context, fileName string, contentType string, expires time.Duration) (objectKey string, url string, err error)
}
