package entity

import (
	"github.com/google/uuid"
)

type ChannelsIntegrationsData struct {
	Code     *string `json:"code,omitempty"`
	Name     *string `json:"name,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
	UserName *string `json:"username,omitempty"`
	Game     *string `json:"game,omitempty"`
	UserId   *string `json:"userId,omitempty"`
}

type DonationAlertsIntegration struct {
	ID            uuid.UUID                 `json:"id"`
	Enabled       bool                      `json:"enabled"`
	ChannelID     string                    `json:"channelId"`
	IntegrationID string                    `json:"integrationId"`
	AccessToken   *string                   `json:"accessToken"`
	RefreshToken  *string                   `json:"refreshToken"`
	ClientID      *string                   `json:"clientId"`
	ClientSecret  *string                   `json:"clientSecret"`
	APIKey        *string                   `json:"apiKey"`
	Data          *ChannelsIntegrationsData `json:"data"`
}
