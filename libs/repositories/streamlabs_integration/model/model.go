package model

import (
	"time"

	"github.com/google/uuid"
)

type StreamlabsIntegration struct {
	ID           uuid.UUID
	Enabled      bool
	ChannelID    string
	AccessToken  string
	RefreshToken string
	UserName     string
	Avatar       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var Nil = StreamlabsIntegration{}
