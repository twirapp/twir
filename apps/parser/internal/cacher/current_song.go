package cacher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/lastfm"
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

	var lastfmService *lastfm.Lastfm
	if lfmIntegration, err := c.services.LastfmRepo.GetByChannelID(
		ctx,
		c.parseCtxChannel.ID,
	); err == nil && !lfmIntegration.IsNil() && lfmIntegration.SessionKey != nil {
		s, lfmErr := lastfm.New(
			lastfm.Opts{
				ApiKey:       c.services.Config.LastFM.ApiKey,
				ClientSecret: c.services.Config.LastFM.ClientSecret,
				SessionKey:   *lfmIntegration.SessionKey,
			},
		)
		if lfmErr != nil {
			c.services.Logger.Error("failed to create lastfm service", zap.Error(lfmErr))
		} else {
			lastfmService = s
		}
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

	var vkService *vk.VK
	vkEntity, err := c.services.VKRepo.GetByChannelID(ctx, c.parseCtxChannel.ID)
	if err == nil && !vkEntity.IsNil() && vkEntity.Enabled && vkEntity.AccessToken != "" {
		v, vkErr := vk.New(
			vk.Opts{
				Integration: vkEntity,
			},
		)
		if vkErr != nil {
			c.services.Logger.Error("failed to create vk service", zap.Error(vkErr))
		} else {
			vkService = v
		}
	}

	integrationsForFetch := []model.IntegrationService{
		model.IntegrationServiceSpotify,
		model.IntegrationServiceVK,
		"YOUTUBE_SR",
		"MUSIC_RECOGNIZER",
		"LASTFM",
	}

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
		case "LASTFM":
			if lastfmService == nil {
				continue
			}

			track, err := lastfmService.GetTrack()
			if err != nil {
				c.services.Logger.Error("failed to get track from lfm", zap.Error(err))
				continue
			}

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
