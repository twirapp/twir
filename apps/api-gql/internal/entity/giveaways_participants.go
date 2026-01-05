package entity

import "github.com/google/uuid"

type ChannelGiveawayParticipant struct {
	UserID      string
	UserLogin   string
	DisplayName string
	IsWinner    bool
	ID          uuid.UUID
	GiveawayID  uuid.UUID
}

var ChannelGiveawayParticipantNil = ChannelGiveawayParticipant{}
