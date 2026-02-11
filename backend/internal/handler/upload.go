package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"upload/internal/entity"
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

type presignRequest struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
}

type presignResponse struct {
	ObjectKey string            `json:"objectKey"`
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	ExpiresIn int               `json:"expiresIn"`
}

func (h *UploadHandler) Presign(w http.ResponseWriter, r *http.Request) {
	var req presignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.FileName == "" || req.ContentType == "" {
		http.Error(w, "fileName and contentType are required", http.StatusBadRequest)
		return
	}
	expires := 15 * time.Minute
	objectKey, url, err := h.presign.PresignPut(r.Context(), req.FileName, req.ContentType, expires)
	if err != nil {
		http.Error(w, "failed to presign", http.StatusInternalServerError)
		return
	}

	resp := presignResponse{
		ObjectKey: objectKey,
		URL:       url,
		Method:    http.MethodPut,
		Headers: map[string]string{
			"Content-Type": req.ContentType,
		},
		ExpiresIn: int(expires.Seconds()),
	}

	writeJSON(w, http.StatusOK, resp)
}

type completeRequest struct {
	ObjectKey   string `json:"objectKey"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}

type completeResponse struct {
	ID        int64  `json:"id"`
	ObjectKey string `json:"objectKey"`
	URL       string `json:"url"`
}

func (h *UploadHandler) Complete(w http.ResponseWriter, r *http.Request) {
	var req completeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.ObjectKey == "" || req.FileName == "" || req.ContentType == "" || req.Size <= 0 {
		http.Error(w, "objectKey, fileName, contentType, size are required", http.StatusBadRequest)
		return
	}

	u := &entity.Upload{
		ObjectKey:   req.ObjectKey,
		FileName:    req.FileName,
		ContentType: req.ContentType,
		Size:        req.Size,
	}
	id, err := h.repo.Create(r.Context(), u)
	if err != nil {
		http.Error(w, "failed to save", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, completeResponse{
		ID:        id,
		ObjectKey: req.ObjectKey,
		URL:       buildPublicURL(h.publicBase, h.bucket, req.ObjectKey),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func buildPublicURL(publicBase string, bucket string, objectKey string) string {
	publicBase = strings.TrimRight(publicBase, "/")
	return publicBase + "/" + bucket + "/" + objectKey
}
