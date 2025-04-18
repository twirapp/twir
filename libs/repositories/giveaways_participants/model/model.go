package model

import "github.com/oklog/ulid/v2"

type ChannelGiveawayParticipant struct {
	DisplayName string    `db:"display_name"`
	UserID      string    `db:"user_id"`
	IsWinner    bool      `db:"is_winner"`
	ID          ulid.ULID `db:"id"`
	GiveawayID  ulid.ULID `db:"giveaway_id"`
}

var ChannelGiveawayParticipantNil = ChannelGiveawayParticipant{}
