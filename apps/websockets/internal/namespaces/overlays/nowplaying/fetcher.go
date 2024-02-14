package nowplaying

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/lastfm"
	"github.com/twirapp/twir/libs/integrations/spotify"
)

func (c *NowPlaying) startTrackUpdater(ctx context.Context, userId string) error {
	ticker := time.NewTicker(1 * time.Second)
	mu := c.redisLock.NewMutex("overlays:nowplaying:lock:" + userId)

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

	redisKey := fmt.Sprintf("overlays:nowplaying:%s", userId)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			_ = mu.Lock()
			track := &Track{}
			err := c.redis.Get(ctx, redisKey).Scan(track)
			if err == nil {
				_ = c.SendEvent(userId, "nowplaying", track)
				continue
			} else if !errors.Is(err, redis.Nil) {
				c.logger.Error("cannot get redis key", slog.Any("err", err))
			}

			track = fetcher.fetch()
			if track != nil {
				err := c.redis.Set(ctx, redisKey, track, 10*time.Second).Err()
				if err != nil {
					c.logger.Error("cannot set redis key", slog.Any("err", err))
				}
			}

			_ = c.SendEvent(userId, "nowplaying", track)
			_, _ = mu.Unlock()
		}
	}
}

type Track struct {
	Artist   string `json:"artist"`
	Title    string `json:"title"`
	ImageUrl string `json:"image_url,omitempty"`
}

func (i Track) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (i *Track) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, i)
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
