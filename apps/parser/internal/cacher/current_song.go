package cacher

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	lfm "github.com/shkh/lastfm-go/lastfm"
	"github.com/twirapp/twir/libs/integrations/spotify"
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
	var vk *vkService
	if ok {
		vk = newVk(vkIntegration)
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
			track := spoti.GetTrack()
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
			if vk == nil {
				continue
			}
			track := vk.GetTrack(ctx)
			if track != nil {
				c.cache.currentSong = &types.CurrentSong{
					Name: *track,
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

type vkService struct {
	integration *model.ChannelsIntegrations
}

func newVk(integration *model.ChannelsIntegrations) *vkService {
	if integration == nil || !integration.AccessToken.Valid {
		return nil
	}

	service := vkService{
		integration: integration,
	}

	return &service
}

type vkError struct {
	Code int    `json:"error_code,omitempty"`
	Msg  string `error_msg:"msg,omitempty"`
}

type vkAudio struct {
	Artist string `json:"artist,omitempty"`
	Title  string `json:"title,omitempty"`
}

type vkStatus struct {
	Text  *string  `json:"text,omitempty"`
	Audio *vkAudio `json:"audio,omitempty"`
}

type vkResponse struct {
	Error  *vkError  `json:"error,omitempty"`
	Status *vkStatus `json:"response"`
}

func (c *vkService) GetTrack(ctx context.Context) *string {
	data := vkResponse{}
	var response string

	resp, err := req.R().
		SetContext(ctx).
		SetQueryParam("access_token", c.integration.AccessToken.String).
		SetQueryParam("v", "5.131").
		SetSuccessResult(&data).
		SetContentType("application/json").
		Get("https://api.vk.com/method/status.get")

	if err != nil || !resp.IsSuccess() {
		return nil
	}

	if data.Error != nil || data.Status == nil || data.Status.Audio == nil {
		return nil
	}

	status := *data.Status.Audio
	response = fmt.Sprintf(
		"%s — %s",
		status.Artist,
		status.Title,
	)

	return &response
}
