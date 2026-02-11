package service

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
)

type PresignService struct {
	client *minio.Client
	bucket string
}

func NewPresignService(client *minio.Client, bucket string) *PresignService {
	return &PresignService{client: client, bucket: bucket}
}

func (s *PresignService) PresignPut(ctx context.Context, fileName string, contentType string, expires time.Duration) (string, string, error) {
	objectKey := buildObjectKey(fileName)
	url, err := s.client.PresignedPutObject(ctx, s.bucket, objectKey, expires)
	if err != nil {
		return "", "", err
	}
	return objectKey, url.String(), nil
}
