package youtube_bot

import (
	"time"

	"github.com/google/uuid"
)

type YouTubeBot struct {
	ID                    uuid.UUID
	EncryptedAccessToken  string
	EncryptedRefreshToken string
	Scopes                []string
	ExpiresIn             int
	ObtainmentTimestamp   time.Time
	YouTubeUserID         uuid.UUID
	CreatedAt             time.Time
	UpdatedAt             time.Time

	isNil bool
}

func (b YouTubeBot) IsNil() bool { return b.isNil }

var Nil = YouTubeBot{isNil: true}
