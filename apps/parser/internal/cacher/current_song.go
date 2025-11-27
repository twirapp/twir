package cacher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	lfm "github.com/shkh/lastfm-go/lastfm"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
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

	var lastfmService *lastFm
	if ok {
		lastfmService = newLastfm(lastFmIntegration)
	}

	var spotifyService *spotify.Spotify
	spotifyEntity, err := c.services.SpotifyRepo.GetByChannelID(ctx, c.parseCtxChannel.ID)
	if err != nil {
		c.services.Logger.Error("failed to get spotify entity", zap.Error(err))
		return nil
	}
	if spotifyEntity.AccessToken != "" {
		spotifyIntegration := model.Integrations{}
		if err := c.services.Gorm.
			Where("service = ?", "SPOTIFY").
			First(&spotifyIntegration).
			Error; err != nil {
			c.services.Logger.Error("failed to get spotify integration", zap.Error(err))
			return nil
		}

		spotifyService = spotify.New(spotifyIntegration, spotifyEntity, c.services.SpotifyRepo)
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

	integrationsForFetch = append(
		integrationsForFetch,
		model.IntegrationServiceSpotify,
		"YOUTUBE_SR",
		"MUSIC_RECOGNIZER",
	)

checkServices:
	for _, integration := range integrationsForFetch {
		switch integration {
		case model.IntegrationServiceSpotify:
			if spotifyService == nil {
				continue
			}
			track, err := spotifyService.GetTrack(ctx)
			if err != nil {
				c.services.Logger.Error("failed to get track", zap.Error(err))
			}
			if track != nil {
				c.cache.currentSong = &types.CurrentSong{
					Name:  track.Artist + " — " + track.Title,
					Image: track.Image,
				}

				if track.Playlist != nil {
					c.cache.currentSong.Playlist = &types.CurrentSongPlayList{
						Href: track.Playlist.ExternalUrl,
					}

					if track.Playlist.Meta != nil {
						c.cache.currentSong.Playlist.Name = &track.Playlist.Meta.Name
						c.cache.currentSong.Playlist.Followers = &track.Playlist.Meta.Followers
						if len(track.Playlist.Meta.Images) > 0 {
							c.cache.currentSong.Playlist.Image = &track.Playlist.Meta.Images[0]
						}
					}
				}

				break checkServices
			}
		case model.IntegrationServiceLastfm:
			if lastfmService == nil {
				continue
			}

			track := lastfmService.GetTrack()

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
		case "MUSIC_RECOGNIZER":
			if c.services.Config.MusicRecognizerAddr == "" {
				continue
			}

			u, _ := url.Parse(c.services.Config.MusicRecognizerAddr)

			query := u.Query()
			query.Set("channel", c.parseCtxChannel.Name)
			u.RawQuery = query.Encode()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
			if err != nil {
				c.services.Logger.Error("failed to create recognize request", zap.Error(err))
				continue
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				c.services.Logger.Error("failed to recognize track", zap.Error(err))
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				c.services.Logger.Error("failed to read recognize response", zap.Error(err))
				continue
			}

			var successResult struct {
				Track struct {
					Title  string `json:"title"`
					Artist string `json:"artist"`
				} `json:"track"`
				Service string `json:"service"`
			}

			if err := json.Unmarshal(body, &successResult); err != nil {
				c.services.Logger.Error("failed to unmarshal recognize response", zap.Error(err))
				continue
			}

			c.cache.currentSong = &types.CurrentSong{
				Name: fmt.Sprintf("%s — %s", successResult.Track.Artist, successResult.Track.Title),
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
