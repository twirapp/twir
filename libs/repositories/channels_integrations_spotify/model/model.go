package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelIntegrationSpotify struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	AccessToken  string
	RefreshToken string
	ChannelID    string
	AvatarURI    string
	Username     string
	Scopes       []string
	ID           uuid.UUID
	Enabled      bool
}

var Nil = ChannelIntegrationSpotify{}
