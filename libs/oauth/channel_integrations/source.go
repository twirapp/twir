package channel_integrations

import (
	"context"
	"errors"
	"fmt"

	"github.com/twirapp/twir/libs/oauth"
	"github.com/twirapp/twir/libs/oauth/nightbot"
	"github.com/twirapp/twir/libs/oauth/spotify"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

var ErrUnsupportedService = errors.New("channel integration OAuth service is unsupported")

type Provider interface {
	Token(context.Context, integrationsmodel.Service, string) (oauth.Credential, error)
}

type Source struct {
	spotify spotify.TokenSource
	nightbot nightbot.TokenSource
}

func New(spotifySource spotify.TokenSource, nightbotSource nightbot.TokenSource) Provider {
	return Source{spotify: spotifySource, nightbot: nightbotSource}
}

func (source Source) Token(ctx context.Context, service integrationsmodel.Service, channelID string) (oauth.Credential, error) {
	switch service {
	case integrationsmodel.ServiceSpotify:
		if source.spotify != nil {
			return source.spotify.Token(ctx, channelID)
		}
	case integrationsmodel.ServiceNightbot:
		if source.nightbot != nil {
			return source.nightbot.Token(ctx, channelID)
		}
	}
	return oauth.Credential{}, fmt.Errorf("%w: %s", ErrUnsupportedService, service)
}
