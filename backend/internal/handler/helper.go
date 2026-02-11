package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func buildPublicURL(publicBase string, bucket string, objectKey string) string {
	publicBase = strings.TrimRight(publicBase, "/")
	return publicBase + "/" + bucket + "/" + objectKey
}
