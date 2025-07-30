package model

import (
	"github.com/google/uuid"
)

type DonatePayIntegration struct {
	ID        uuid.UUID
	ChannelID string
	ApiKey    string
}
