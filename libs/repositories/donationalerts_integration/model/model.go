package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type DonationAlertsIntegration struct {
	ID           int64
	PublicID     ulid.ULID
	Enabled      bool
	ChannelID    string
	AccessToken  string
	RefreshToken string
	UserName     string
	Avatar       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var Nil = DonationAlertsIntegration{}
