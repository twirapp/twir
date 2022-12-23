package lastfm

import (
	"fmt"

	model "github.com/satont/tsuwari/libs/gomodels"

	lfm "github.com/shkh/lastfm-go/lastfm"
)

type LastFm struct {
	integration *model.ChannelsIntegrations
}

func New(integration *model.ChannelsIntegrations) *LastFm {
	if integration == nil || !integration.APIKey.Valid || !integration.Integration.APIKey.Valid ||
		!integration.Integration.ClientSecret.Valid {
		return nil
	}

	service := LastFm{
		integration: integration,
	}

	return &service
}

func (c *LastFm) GetTrack() *string {
	api := lfm.New(
		c.integration.Integration.APIKey.String,
		c.integration.Integration.ClientSecret.String,
	)
	api.SetSession(c.integration.APIKey.String)

	user, err := api.User.GetInfo(map[string]interface{}{})
	if err != nil {
		return nil
	}

	tracks, err := api.User.GetRecentTracks(map[string]interface{}{
		"limit": "1",
		"user":  user.Name,
	})

	if err != nil || len(tracks.Tracks) == 0 || tracks.Tracks[0].NowPlaying != "true" {
		return nil
	}

	track := tracks.Tracks[0]

	response := fmt.Sprintf("%s — %s", track.Artist.Name, track.Name)

	return &response
}
