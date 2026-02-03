package model

import (
	"time"
)

type Pastebin struct {
	ID          string
	CreatedAt   time.Time
	Content     string
	ExpireAt    *time.Time
	OwnerUserID *string
	UserIP      *string
	UserAgent   *string
}

var Nil = Pastebin{}
