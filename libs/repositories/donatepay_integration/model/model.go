package model

import (
	"github.com/oklog/ulid/v2"
)

type DonatePayIntegration struct {
	ID        ulid.ULID
	ChannelID string
	ApiKey    string
	Enabled   bool
}
