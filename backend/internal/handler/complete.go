package handler

import (
	"encoding/json"
	"net/http"

	"upload/internal/entity"
)

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
