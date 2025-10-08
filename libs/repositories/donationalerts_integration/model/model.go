package model

import (
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type DonationAlertsIntegration struct {
	ID            uuid.UUID
	Enabled       bool
	ChannelID     string
	IntegrationID string
	AccessToken   *string
	RefreshToken  *string
	ClientID      *string
	ClientSecret  *string
	APIKey        *string
	Data          json.RawMessage // JSONB data
}

var Nil = DonationAlertsIntegration{}
