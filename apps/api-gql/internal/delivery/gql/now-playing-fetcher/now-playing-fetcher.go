package now_playing_fetcher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/goccy/go-json"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	buscore "github.com/twirapp/twir/libs/bus-core"
	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/lastfm"
	"github.com/twirapp/twir/libs/integrations/spotify"
	"github.com/twirapp/twir/libs/integrations/vk"
	"github.com/twirapp/twir/libs/logger"
	channelsintegrationslastfm "github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
	vkintegration "github.com/twirapp/twir/libs/repositories/vk_integration"
)

const spotifyTokenRefreshCooldown = 5 * time.Minute

type Opts struct {
	Logger            *slog.Logger
	SpotifyRepository channelsintegrationsspotify.Repository
	LastfmRepository  channelsintegrationslastfm.Repository
	VKRepository      vkintegration.Repository
	Config            cfg.Config
	TwirBus           *buscore.Bus
	Kv                kv.KV
	ChannelID         string
}

type NowPlayingFetcher struct {
	spotifyRepository channelsintegrationsspotify.Repository
	logger            *slog.Logger
	kv                kv.KV
	twirBus           *buscore.Bus

	channelId                 string
	spotifyScopes             []string
	lastSpotifyTokenRefreshAt time.Time

	lastfmService  *lastfm.Lastfm
	spotifyService *spotify.Spotify
	vkService      *vk.VK
}

func New(opts Opts) (*NowPlayingFetcher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var lfmService *lastfm.Lastfm
	var vkService *vk.VK

	// Get lastfm integration from the new repository
	lastfmIntegration, err := opts.LastfmRepository.GetByChannelID(ctx, opts.ChannelID)
	if err == nil && !lastfmIntegration.IsNil() && lastfmIntegration.Enabled {
		var sessionKey, userName string
		if lastfmIntegration.SessionKey != nil {
			sessionKey = *lastfmIntegration.SessionKey
		}
		if lastfmIntegration.UserName != nil {
			userName = *lastfmIntegration.UserName
		}

		l, err := lastfm.New(
			ctx,
			lastfm.Opts{
				ApiKey:       opts.Config.LastFM.ApiKey,
				ClientSecret: opts.Config.LastFM.ClientSecret,
				SessionKey:   sessionKey,
				UserName:     userName,
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
	spotifyEnabled := spotifyEntity.AccessToken != "" && spotifyEntity.RefreshToken != ""

	// Get VK integration from the new repository
	vkEntity, err := opts.VKRepository.GetByChannelID(ctx, opts.ChannelID)
	if err == nil && !vkEntity.IsNil() && vkEntity.Enabled && vkEntity.AccessToken != "" {
		v, err := vk.New(
			vk.Opts{
				Integration: vkEntity,
			},
		)
		if err == nil {
			vkService = v
		}
	}

	f := &NowPlayingFetcher{
		spotifyRepository: opts.SpotifyRepository,
		channelId:         opts.ChannelID,
		kv:                opts.Kv,
		lastfmService:     lfmService,
		vkService:         vkService,
		logger:            opts.Logger,
		twirBus:           opts.TwirBus,
		spotifyScopes:     spotifyEntity.Scopes,
	}

	if spotifyEnabled {
		f.maybeRefreshSpotifyService(ctx)
	}

	return f, nil
}

func (c *NowPlayingFetcher) maybeRefreshSpotifyService(ctx context.Context) {
	if c.twirBus == nil || c.channelId == "" || len(c.spotifyScopes) == 0 {
		return
	}
	if !c.lastSpotifyTokenRefreshAt.IsZero() && time.Since(c.lastSpotifyTokenRefreshAt) < spotifyTokenRefreshCooldown {
		return
	}
	c.lastSpotifyTokenRefreshAt = time.Now()

	token, err := c.twirBus.Tokens.RequestChannelIntegrationToken.Request(
		ctx,
		buscoretokens.GetChannelIntegrationTokenRequest{
			ChannelID: c.channelId,
			Service:   integrationsmodel.ServiceSpotify,
		},
	)
	if err != nil {
		c.logger.Error(
			"failed to get spotify token from tokens service",
			logger.Error(err),
			slog.String("channel_id", c.channelId),
		)
		return
	}

	c.spotifyService = spotify.NewStatic(token.Data.AccessToken, c.spotifyScopes)
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
		cachedTrack.advanceProgress(time.Now())
		cachedTrack.fromCache = true
		return cachedTrack, nil
	} else if !errors.Is(err, kv.ErrKeyNil) {
		return nil, err
	}

	if c.spotifyService == nil {
		c.maybeRefreshSpotifyService(ctx)
	}

	if c.spotifyService != nil {
		spotifyTrack, err := c.spotifyService.GetTrack(ctx)
		if err != nil {
			c.maybeRefreshSpotifyService(ctx)
			if c.spotifyService != nil {
				spotifyTrack, err = c.spotifyService.GetTrack(ctx)
			}
		}
		if err != nil {
			c.logger.Error(
				"cannot fetch spotify track",
				logger.Error(err),
				slog.String("channel_id", c.channelId),
			)
		}

		if spotifyTrack != nil && spotifyTrack.IsPlaying {
			progressMs := spotifyTrack.ProgressMs
			durationMs := spotifyTrack.DurationMs
			progressObservedAt := time.Now()

			return &Track{
				Artist:             spotifyTrack.Artist,
				Title:              spotifyTrack.Title,
				ImageUrl:           spotifyTrack.Image,
				ProgressMs:         &progressMs,
				DurationMs:         &durationMs,
				ProgressObservedAt: progressObservedAt,
			}, nil
		}
	}

	if c.lastfmService != nil {
		lastfmTrack, err := c.lastfmService.GetTrack(ctx)
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
	Artist             string    `json:"artist"`
	Title              string    `json:"title"`
	ImageUrl           string    `json:"image_url,omitempty"`
	ProgressMs         *int      `json:"progress_ms,omitempty"`
	DurationMs         *int      `json:"duration_ms,omitempty"`
	ProgressObservedAt time.Time `json:"progress_observed_at,omitempty"`
	fromCache          bool
}

func (t *Track) advanceProgress(now time.Time) {
	if t.ProgressMs == nil || t.DurationMs == nil || *t.DurationMs <= 0 || t.ProgressObservedAt.IsZero() {
		return
	}

	elapsed := now.Sub(t.ProgressObservedAt)
	if elapsed < 0 {
		elapsed = 0
	}

	progress := *t.ProgressMs + int(elapsed/time.Millisecond)
	if progress < 0 {
		progress = 0
	}
	if progress > *t.DurationMs {
		progress = *t.DurationMs
	}

	t.ProgressMs = &progress
	t.ProgressObservedAt = now
}

func (i Track) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (i *Track) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, i)
}
