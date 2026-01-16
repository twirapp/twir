package entity

import (
	"time"
)

type ShortenedUrl struct {
	ID          string
	Link        string
	Views       int
	OwnerUserID *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var Nil = ShortenedUrl{}
