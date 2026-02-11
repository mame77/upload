package service

import (
	"path"
	"strings"

	"github.com/google/uuid"
)

func buildObjectKey(fileName string) string {
	ext := strings.ToLower(path.Ext(fileName))
	if ext == "" || len(ext) > 16 {
		ext = ""
	}
	return uuid.New().String() + ext
}
