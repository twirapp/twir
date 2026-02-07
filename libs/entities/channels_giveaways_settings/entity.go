package channels_giveaways_settings

import (
	"time"

	"github.com/google/uuid"
)

type Settings struct {
	ID            uuid.UUID
	ChannelID     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	WinnerMessage string

	isNil bool
}

func (s Settings) IsNil() bool {
	return s.isNil
}

var Nil = Settings{
	isNil: true,
}
