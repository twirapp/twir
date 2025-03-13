package model

import (
	"github.com/google/uuid"
)

type SevenTvIntegration struct {
	ID                         uuid.UUID
	RewardIdForAddEmote        *string
	RewardIdForRemoveEmote     *string
	DeleteEmotesOnlyAddedByApp bool
	AddedEmotes                []string
	ChannelID                  string
}

var Nil = SevenTvIntegration{}
