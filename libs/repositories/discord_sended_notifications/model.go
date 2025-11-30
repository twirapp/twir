package discord_sended_notifications

import (
	"time"
)

type DiscordSendedNotification struct {
	ID               string    `db:"id"`
	GuildID          string    `db:"guild_id"`
	MessageID        string    `db:"message_id"`
	TwitchChannelID  string    `db:"channel_id"`
	DiscordChannelID string    `db:"discord_channel_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`

	isNil bool
}

func (d DiscordSendedNotification) IsNil() bool {
	return d.isNil
}

var Nil = DiscordSendedNotification{isNil: true}
