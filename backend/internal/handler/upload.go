package handler

import (
	"upload/internal/port"
)

type UploadHandler struct {
	repo       port.UploadRepository
	presign    port.PresignService
	publicBase string
	bucket     string
}

func NewUploadHandler(repo port.UploadRepository, presign port.PresignService, publicBase string, bucket string) *UploadHandler {
	return &UploadHandler{repo: repo, presign: presign, publicBase: publicBase, bucket: bucket}
}
