package model

import "github.com/oklog/ulid/v2"

type ChannelGiveawayParticipant struct {
	UserID      string    `db:"user_id"`
	UserLogin   string    `db:"user_login"`
	DisplayName string    `db:"display_name"`
	IsWinner    bool      `db:"is_winner"`
	ID          ulid.ULID `db:"id"`
	GiveawayID  ulid.ULID `db:"giveaway_id"`
}

var ChannelGiveawayParticipantNil = ChannelGiveawayParticipant{}
