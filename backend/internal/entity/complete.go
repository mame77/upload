package entity

import "time"

type Upload struct {
	ID          int64
	ObjectKey   string
	FileName    string
	ContentType string
	Size        int64
	CreatedAt   time.Time
}
