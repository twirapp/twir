package entity

import (
	"time"
)

type Pastebin struct {
	ID          string     `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	Content     string     `json:"content"`
	ExpireAt    *time.Time `json:"expire_at"`
	OwnerUserID *string    `json:"owner_user_id"`
}

var PastebinNil = Pastebin{}
