package now_playing_fetcher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	cfg "github.com/twirapp/twir/libs/config"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/lastfm"
	"github.com/twirapp/twir/libs/integrations/spotify"
	"github.com/twirapp/twir/libs/integrations/vk"
	"github.com/twirapp/twir/libs/logger"
	channelsintegrationslastfm "github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	"gorm.io/gorm"
)

type Opts struct {
	Logger            *slog.Logger
	SpotifyRepository channelsintegrationsspotify.Repository
	LastfmRepository  channelsintegrationslastfm.Repository
	Config            cfg.Config
	Gorm              *gorm.DB
	Kv                kv.KV
	ChannelID         string
}

type NowPlayingFetcher struct {
	spotifyRepository channelsintegrationsspotify.Repository
	logger            *slog.Logger

	gorm *gorm.DB
	kv   kv.KV

	lastfmService  *lastfm.Lastfm
	spotifyService *spotify.Spotify
	vkService      *vk.VK
	channelId      string
}

func New(opts Opts) (*NowPlayingFetcher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var channelIntegrations []*model.ChannelsIntegrations
	if err := opts.Gorm.
		Where(`"channelId" = ?`, opts.ChannelID).
		Preload("Integration").
		Find(&channelIntegrations).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get channel integrations: %w", err)
	}

	vkEntity, _ := lo.Find(
		channelIntegrations,
		func(integration *model.ChannelsIntegrations) bool {
			return integration.Integration.Service == "VK" && integration.Enabled
		},
	)

	var lfmService *lastfm.Lastfm
	var spotifyService *spotify.Spotify
	var vkService *vk.VK

	// Get lastfm integration from the new repository
	lastfmIntegration, err := opts.LastfmRepository.GetByChannelID(ctx, opts.ChannelID)
	if err == nil && !lastfmIntegration.IsNil() && lastfmIntegration.Enabled && lastfmIntegration.SessionKey != nil {
		l, err := lastfm.New(
			lastfm.Opts{
				ApiKey:       opts.Config.LastFM.ApiKey,
				ClientSecret: opts.Config.LastFM.ClientSecret,
				SessionKey:   *lastfmIntegration.SessionKey,
			},
		)
		if err == nil {
			lfmService = l
		}
	}

	spotifyEntity, err := opts.SpotifyRepository.GetByChannelID(ctx, opts.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spotify integration: %w", err)
	}
	if spotifyEntity.AccessToken != "" && spotifyEntity.RefreshToken != "" {
		spotifyIntegration := model.Integrations{}
		if err := opts.Gorm.
			Where("service = ?", "SPOTIFY").
			First(&spotifyIntegration).
			Error; err != nil {
			return nil, fmt.Errorf("failed to get spotify integration: %w", err)
		}

		spotifyService = spotify.New(spotifyIntegration, spotifyEntity, opts.SpotifyRepository)
	}

	if vkEntity != nil {
		v, err := vk.New(
			vk.Opts{
				Gorm:        opts.Gorm,
				Integration: vkEntity,
			},
		)
		if err == nil {
			vkService = v
		}
	}

	return &NowPlayingFetcher{
		spotifyRepository: opts.SpotifyRepository,
		channelId:         opts.ChannelID,
		gorm:              opts.Gorm,
		kv:                opts.Kv,
		lastfmService:     lfmService,
		spotifyService:    spotifyService,
		vkService:         vkService,
		logger:            opts.Logger,
	}, nil
}

func (c *NowPlayingFetcher) Fetch(ctx context.Context) (*Track, error) {
	track, err := c.fetchWrapper(ctx)
	if err != nil {
		return nil, err
	}

	if track != nil && !track.fromCache {
		redisKey := fmt.Sprintf("overlays:nowplaying:%s", c.channelId)
		if err := c.kv.Set(
			ctx,
			redisKey,
			track,
			kvoptions.WithExpire(10*time.Second),
		); err != nil {
			return nil, err
		}
	}

	return track, nil
}

func (c *NowPlayingFetcher) fetchWrapper(ctx context.Context) (*Track, error) {
	redisKey := fmt.Sprintf("overlays:nowplaying:%s", c.channelId)

	cachedTrack := &Track{}
	err := c.kv.Get(ctx, redisKey).Scan(cachedTrack)
	if err == nil {
		cachedTrack.fromCache = true
		return cachedTrack, nil
	} else if !errors.Is(err, kv.ErrKeyNil) {
		return nil, err
	}

	if c.spotifyService != nil {
		spotifyTrack, err := c.spotifyService.GetTrack(ctx)
		if err != nil {
			c.logger.Error(
				"cannot fetch spotify track",
				logger.Error(err),
				slog.String("channel_id", c.channelId),
			)
		}

		if spotifyTrack != nil && spotifyTrack.IsPlaying {
			return &Track{
				Artist:   spotifyTrack.Artist,
				Title:    spotifyTrack.Title,
				ImageUrl: spotifyTrack.Image,
			}, nil
		}
	}

	if c.lastfmService != nil {
		lastfmTrack, err := c.lastfmService.GetTrack()
		if err != nil {
			c.logger.Error(
				"cannot fetch lastfm track",
				logger.Error(err),
				slog.String("channel_id", c.channelId),
			)
		}

		if lastfmTrack != nil {
			return &Track{
				Artist:   lastfmTrack.Artist,
				Title:    lastfmTrack.Title,
				ImageUrl: lastfmTrack.Image,
			}, nil
		}
	}

	if c.vkService != nil {
		vkTrack, err := c.vkService.GetTrack(ctx)
		if err != nil {
			c.logger.Error(
				"cannot fetch vk track",
				logger.Error(err),
				slog.String("channel_id", c.channelId),
			)
		}

		if vkTrack != nil {
			return &Track{
				Artist:   vkTrack.Artist,
				Title:    vkTrack.Title,
				ImageUrl: vkTrack.Image,
			}, nil
		}
	}

	return nil, nil
}

type Track struct {
	Artist    string `json:"artist"`
	Title     string `json:"title"`
	ImageUrl  string `json:"image_url,omitempty"`
	fromCache bool   `json:"from_cache,omitempty"`
}

func (i Track) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (i *Track) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, i)
}
