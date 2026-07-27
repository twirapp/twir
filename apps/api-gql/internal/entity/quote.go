package entity

import (
	"time"

	"github.com/google/uuid"
)

type Quote struct {
	ID          uuid.UUID
	ChannelID   string
	Number      int
	Text        string
	CreatorID   *string
	CreatorName *string
	GameID      *string
	GameName    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var QuoteNil = Quote{}
