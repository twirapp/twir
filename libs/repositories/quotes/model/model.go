package model

import (
	"time"

	"github.com/google/uuid"
)

type Quote struct {
	ID          uuid.UUID `db:"id"`
	ChannelID   uuid.UUID `db:"channel_id"`
	Number      int       `db:"number"`
	Text        string    `db:"text"`
	CreatorID   *string   `db:"creator_id"`
	CreatorName *string   `db:"creator_name"`
	GameID      *string   `db:"game_id"`
	GameName    *string   `db:"game_name"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

var Nil = Quote{}
