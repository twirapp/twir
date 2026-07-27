package model

import (
	"time"

	"github.com/google/uuid"
)

type Quote struct {
	ID          uuid.UUID
	ChannelID   uuid.UUID
	Number      int
	Text        string
	CreatorID   *string
	CreatorName *string
	GameID      *string
	GameName    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var Nil = Quote{}
