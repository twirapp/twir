package model

import (
	"time"
)

type ShortenedUrl struct {
	ShortID         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	URL             string
	CreatedByUserId *string
	Views           int
}

var Nil = ShortenedUrl{}
