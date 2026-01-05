package model

import "github.com/google/uuid"

type ChannelGiveawayParticipant struct {
	UserID      string    `db:"user_id"`
	UserLogin   string    `db:"user_login"`
	DisplayName string    `db:"display_name"`
	IsWinner    bool      `db:"is_winner"`
	ID          uuid.UUID `db:"id"`
	GiveawayID  uuid.UUID `db:"giveaway_id"`
}

var ChannelGiveawayParticipantNil = ChannelGiveawayParticipant{}
