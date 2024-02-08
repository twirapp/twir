package nowplaying

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/lastfm"
	"github.com/twirapp/twir/libs/integrations/spotify"
)

func (c *NowPlaying) fetcher(ctx context.Context, userId string) error {
	ticker := time.NewTicker(1 * time.Second)

	var channelIntegrations []*model.ChannelsIntegrations
	if err := c.gorm.
		WithContext(ctx).
		Where(`"channelId" = ?`, userId).
		Preload("Integration").
		Find(&channelIntegrations).
		Error; err != nil {
		return err
	}

	lfmEntity, _ := lo.Find(
		channelIntegrations,
		func(integration *model.ChannelsIntegrations) bool {
			return integration.Integration.Service == "LASTFM" && integration.Enabled
		},
	)
	spotifyEntity, _ := lo.Find(
		channelIntegrations,
		func(integration *model.ChannelsIntegrations) bool {
			return integration.Integration.Service == "SPOTIFY" && integration.Enabled
		},
	)

	var lfm *lastfm.Lastfm
	var spoti *spotify.Spotify

	if lfmEntity != nil {
		l, err := lastfm.New(
			lastfm.Opts{
				Gorm:        c.gorm,
				Integration: lfmEntity,
			},
		)
		if err != nil {
			return err
		}

		lfm = l
	}

	if spotifyEntity != nil {
		spoti = spotify.New(spotifyEntity, c.gorm)
	}

	fetcher := &channelSongFetcher{
		lfm:   lfm,
		spoti: spoti,
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			cachedValue := c.redis.Get(ctx, fmt.Sprintf("overlays:nowplaying:%s", userId)).Val()
			if cachedValue != "" {
				_ = c.SendEvent(userId, "nowplaying", cachedValue)
				continue
			}

			track := fetcher.fetch()
			if track != nil {
				c.redis.Set(ctx, fmt.Sprintf("overlays:nowplaying:%s", userId), track, 10*time.Second)
			}

			_ = c.SendEvent(userId, "nowplaying", track)
		}
	}
}

type Track struct {
	Artist   string `json:"artist"`
	Title    string `json:"title"`
	ImageUrl string `json:"image_url,omitempty"`
}

type channelSongFetcher struct {
	lfm   *lastfm.Lastfm
	spoti *spotify.Spotify
}

func (c *channelSongFetcher) fetch() *Track {
	if c.spoti != nil {
		spotifyTrack := c.spoti.GetTrack()
		if spotifyTrack != nil {
			return &Track{
				Artist:   spotifyTrack.Artist,
				Title:    spotifyTrack.Title,
				ImageUrl: spotifyTrack.Image,
			}
		}
	}

	if c.lfm != nil {
		lastfmTrack, err := c.lfm.GetTrack()
		if err != nil {
			return nil
		}
		if lastfmTrack != nil {
			return &Track{
				Artist:   lastfmTrack.Artist,
				Title:    lastfmTrack.Title,
				ImageUrl: lastfmTrack.Image,
			}
		}
	}

	return nil
}
