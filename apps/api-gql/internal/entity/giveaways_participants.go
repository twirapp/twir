package entity

import "github.com/oklog/ulid/v2"

type ChannelGiveawayParticipant struct {
	DisplayName string
	UserID      string
	IsWinner    bool
	ID          ulid.ULID
	GiveawayID  ulid.ULID
}

var ChannelGiveawayParticipantNil = ChannelGiveawayParticipant{}
