package cacher

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	lfm "github.com/shkh/lastfm-go/lastfm"
	"github.com/twirapp/twir/libs/integrations/spotify"
	"github.com/twirapp/twir/libs/integrations/vk"
	"go.uber.org/zap"
)

func (c *cacher) GetCurrentSong(ctx context.Context) *types.CurrentSong {
	c.locks.currentSong.Lock()
	defer c.locks.currentSong.Unlock()

	if c.cache.currentSong != nil {
		return c.cache.currentSong
	}

	integrations := c.GetEnabledChannelIntegrations(ctx)
	integrations = lo.Filter(
		integrations,
		func(integration *model.ChannelsIntegrations, _ int) bool {
			switch integration.Integration.Service {
			case "SPOTIFY", "VK", "LASTFM":
				return integration.Enabled
			default:
				return false
			}
		},
	)

	lastFmIntegration, ok := lo.Find(
		integrations,
		func(integration *model.ChannelsIntegrations) bool {
			return integration.Integration.Service == "LASTFM"
		},
	)

	var lfm *lastFm
	if ok {
		lfm = newLastfm(lastFmIntegration)
	}

	spotifyIntegration, ok := lo.Find(
		integrations,
		func(integration *model.ChannelsIntegrations) bool {
			return integration.Integration.Service == "SPOTIFY"
		},
	)
	var spoti *spotify.Spotify
	if ok {
		spoti = spotify.New(spotifyIntegration, c.services.Gorm)
	}

	vkIntegration, ok := lo.Find(
		integrations,
		func(integration *model.ChannelsIntegrations) bool {
			return integration.Integration.Service == "VK"
		},
	)
	var vkService *vk.VK
	if ok {
		vkService, _ = vk.New(
			vk.Opts{
				Gorm:        c.services.Gorm,
				Integration: vkIntegration,
			},
		)
	}

	integrationsForFetch := lo.Map(
		integrations,
		func(integration *model.ChannelsIntegrations, _ int) model.IntegrationService {
			return integration.Integration.Service
		},
	)

	integrationsForFetch = append(integrationsForFetch, model.IntegrationService("YOUTUBE_SR"))

checkServices:
	for _, integration := range integrationsForFetch {
		switch integration {
		case model.IntegrationServiceSpotify:
			if spoti == nil {
				continue
			}
			track, err := spoti.GetTrack()
			if err != nil {
				c.services.Logger.Error("failed to get track", zap.Error(err))
			}
			if track != nil {
				c.cache.currentSong = &types.CurrentSong{
					Name:  track.Artist + " — " + track.Title,
					Image: track.Image,
				}
				break checkServices
			}
		case model.IntegrationServiceLastfm:
			if lfm == nil {
				continue
			}

			track := lfm.GetTrack()

			if track != nil {
				c.cache.currentSong = &types.CurrentSong{
					Name:  fmt.Sprintf("%s — %s", track.Artist, track.Title),
					Image: track.Image,
				}
				break checkServices
			}
		case model.IntegrationServiceVK:
			if vkService == nil {
				continue
			}
			track, err := vkService.GetTrack(ctx)
			if err != nil {
				c.services.Logger.Error("failed to get track", zap.Error(err))
			}
			if track != nil {
				c.cache.currentSong = &types.CurrentSong{
					Name: fmt.Sprintf("%s — %s", track.Artist, track.Title),
				}
				break checkServices
			}
		case "YOUTUBE_SR":
			redisData, err := c.services.Redis.Get(
				context.Background(),
				fmt.Sprintf("songrequests:youtube:%s:currentPlaying", c.parseCtxChannel.ID),
			).Result()
			if err == redis.Nil {
				continue
			}
			if err != nil {
				continue
			}
			song := model.RequestedSong{}
			if err = c.services.Gorm.
				WithContext(ctx).
				Where("id = ?", redisData).
				First(&song).Error; err != nil {
				fmt.Println("song nog found", err)
				continue
			}

			c.cache.currentSong = &types.CurrentSong{
				Name: fmt.Sprintf(
					`"%s" youtu.be/%s requested by @%s`,
					song.Title,
					song.VideoID,
					song.OrderedByName,
				),
			}
			break checkServices
		}
	}

	return c.cache.currentSong
}

type lastFm struct {
	integration *model.ChannelsIntegrations
}

func newLastfm(integration *model.ChannelsIntegrations) *lastFm {
	if integration == nil || !integration.APIKey.Valid || !integration.Integration.APIKey.Valid ||
		!integration.Integration.ClientSecret.Valid {
		return nil
	}

	service := lastFm{
		integration: integration,
	}

	return &service
}

type LFMGetTrackResponse struct {
	Title  string
	Artist string
	Image  string
}

func (c *lastFm) GetTrack() *LFMGetTrackResponse {
	api := lfm.New(
		c.integration.Integration.APIKey.String,
		c.integration.Integration.ClientSecret.String,
	)
	api.SetSession(c.integration.APIKey.String)

	user, err := api.User.GetInfo(map[string]interface{}{})
	if err != nil {
		return nil
	}

	tracks, err := api.User.GetRecentTracks(
		map[string]interface{}{
			"limit": "1",
			"user":  user.Name,
		},
	)

	if err != nil || len(tracks.Tracks) == 0 || tracks.Tracks[0].NowPlaying != "true" {
		return nil
	}

	track := tracks.Tracks[0]

	// track.Images
	var cover string
	if len(track.Images) > 0 {
		cover = track.Images[0].Url
	}

	return &LFMGetTrackResponse{
		Title:  track.Name,
		Artist: track.Artist.Name,
		Image:  cover,
	}
}
