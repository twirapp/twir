package now_playing_fetcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/kr/pretty"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/lastfm"
	"github.com/twirapp/twir/libs/integrations/spotify"
	"github.com/twirapp/twir/libs/integrations/vk"
	"gorm.io/gorm"
)

type Opts struct {
	Gorm      *gorm.DB
	ChannelID string
	Redis     *redis.Client
}

type NowPlayingFetcher struct {
	channelId string
	gorm      *gorm.DB
	redis     *redis.Client

	lastfmService  *lastfm.Lastfm
	spotifyService *spotify.Spotify
	vkService      *vk.VK
}

func New(opts Opts) (*NowPlayingFetcher, error) {
	var channelIntegrations []*model.ChannelsIntegrations
	if err := opts.Gorm.
		Where(`"channelId" = ?`, opts.ChannelID).
		Preload("Integration").
		Find(&channelIntegrations).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get channel integrations: %w", err)
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
	vkEntity, _ := lo.Find(
		channelIntegrations,
		func(integration *model.ChannelsIntegrations) bool {
			return integration.Integration.Service == "VK" && integration.Enabled
		},
	)

	var lfmService *lastfm.Lastfm
	var spotifyService *spotify.Spotify
	var vkService *vk.VK

	if lfmEntity != nil {
		l, err := lastfm.New(
			lastfm.Opts{
				Gorm:        opts.Gorm,
				Integration: lfmEntity,
			},
		)
		if err == nil {
			lfmService = l
		}
	}

	if spotifyEntity != nil {
		spotifyService = spotify.New(spotifyEntity, opts.Gorm)
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
		channelId:      opts.ChannelID,
		gorm:           opts.Gorm,
		redis:          opts.Redis,
		lastfmService:  lfmService,
		spotifyService: spotifyService,
		vkService:      vkService,
	}, nil
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

func (c *NowPlayingFetcher) Fetch(ctx context.Context) (*Track, error) {
	track, err := c.fetchWrapper(ctx)
	if c.channelId == "115133555" {
		pretty.Println("track for suhodolskiy", track, err)
	}
	if err != nil {
		return nil, err
	}

	if track != nil && !track.fromCache {
		redisKey := fmt.Sprintf("overlays:nowplaying:%s", c.channelId)
		if err := c.redis.Set(ctx, redisKey, track, 10*time.Second).Err(); err != nil {
			return nil, err
		}
	}

	return track, nil
}

func (c *NowPlayingFetcher) fetchWrapper(ctx context.Context) (*Track, error) {
	redisKey := fmt.Sprintf("overlays:nowplaying:%s", c.channelId)

	cachedTrack := &Track{}
	err := c.redis.Get(ctx, redisKey).Scan(cachedTrack)
	if err == nil {
		cachedTrack.fromCache = true
		return cachedTrack, nil
	} else if !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if c.spotifyService != nil {
		spotifyTrack := c.spotifyService.GetTrack()
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
			return nil, err
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
			return nil, err
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
