package entity

import (
	"time"

	"github.com/google/uuid"
)

type DonationAlertsIntegration struct {
	ID           int64
	PublicID     uuid.UUID
	Enabled      bool
	ChannelID    string
	AccessToken  string
	RefreshToken string
	UserName     string
	Avatar       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
