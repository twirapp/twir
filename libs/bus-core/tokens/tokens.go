package tokens

import (
	"github.com/google/uuid"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

const (
	RequestAppTokenSubject                = "tokens.request_app_token"
	RequestUserTokenSubject               = "tokens.request_user_token"
	RequestBotTokenSubject                = "tokens.request_bot_token"
	RequestChannelIntegrationTokenSubject = "tokens.request_channel_integration_token"
)

type GetUserTokenRequest struct {
	UserId uuid.UUID `json:"userId"`
}

type GetBotTokenRequest struct {
	BotId    string                  `json:"botId"`
	Platform platformentity.Platform `json:"platform"`
}

type GetAppTokenRequest struct {
	Platform platformentity.Platform `json:"platform"`
}

type GetChannelIntegrationTokenRequest struct {
	ChannelID string                    `json:"channelId"`
	Service   integrationsmodel.Service `json:"service"`
}

type UpdateTokenRequest struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	ExpiresIn    int64    `json:"expiresIn"`
	Scopes       []string `json:"scopes"`
}

type TokenResponse struct {
	AccessToken string   `json:"access_token"`
	Scopes      []string `json:"scopes"`
	ExpiresIn   int32    `json:"expires_in"`
}
