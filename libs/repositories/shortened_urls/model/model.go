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
}

var Nil = ShortenedUrl{}
