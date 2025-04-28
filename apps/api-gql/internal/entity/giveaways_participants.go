package entity

import "github.com/oklog/ulid/v2"

type ChannelGiveawayParticipant struct {
	UserID      string
	UserLogin   string
	DisplayName string
	IsWinner    bool
	ID          ulid.ULID
	GiveawayID  ulid.ULID
}

var ChannelGiveawayParticipantNil = ChannelGiveawayParticipant{}
