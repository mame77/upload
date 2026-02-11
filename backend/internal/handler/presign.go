package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

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
