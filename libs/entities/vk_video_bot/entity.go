package vk_video_bot

import (
	"time"

	"github.com/google/uuid"
)

type VKVideoBot struct {
	ID                    uuid.UUID
	EncryptedAccessToken  string
	EncryptedRefreshToken string
	Scopes                []string
	ExpiresIn             int
	ObtainmentTimestamp   time.Time
	VKUserID              uuid.UUID
	CreatedAt             time.Time
	UpdatedAt             time.Time

	isNil bool
}

func (b VKVideoBot) IsNil() bool { return b.isNil }

var Nil = VKVideoBot{isNil: true}
