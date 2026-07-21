package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
)

type ChannelPlatform struct {
	ID                uuid.UUID         `db:"id"`
	ChannelID         uuid.UUID         `db:"channel_id"`
	Platform          platform.Platform `db:"platform"`
	UserID            uuid.UUID         `db:"user_id"`
	PlatformChannelID string            `db:"platform_channel_id"`
	Enabled           bool              `db:"enabled"`
	BotUserID         *uuid.UUID        `db:"bot_user_id"`
	BotConfig         json.RawMessage   `db:"bot_config"`
	CreatedAt         time.Time         `db:"created_at"`
	UpdatedAt         time.Time         `db:"updated_at"`

	isNil bool
}

func (c ChannelPlatform) IsNil() bool {
	return c.isNil
}

var Nil = ChannelPlatform{isNil: true}
