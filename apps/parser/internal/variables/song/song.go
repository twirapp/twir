package song

import (
	"context"
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"

	lastfm "github.com/satont/tsuwari/apps/parser/internal/integrations/lastfm"
	vkIntegr "github.com/satont/tsuwari/apps/parser/internal/integrations/vk"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/go-redis/redis/v9"
	spotify "github.com/satont/tsuwari/libs/integrations/spotify"

	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
)

const (
	SOUNDTRACK = "TWITCH_SOUNDTRACK"
	VK         = "VK"
	SPOTIFY    = "SPOTIFY"
	LASTFM     = "LASTFM"
	YOUTUBE_SR = "YOUTUBE_SR"
)

var Variable = types.Variable{
	Name:        "currentsong",
	Description: lo.ToPtr("Current played song"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		db := do.MustInvoke[gorm.DB](di.Provider)
		redisClient := do.MustInvoke[redis.Client](di.Provider)

		result := &types.VariableHandlerResult{}

		integrations := ctx.GetEnabledIntegrations()

		integrations = lo.Filter(
			integrations,
			func(integration model.ChannelsIntegrations, _ int) bool {
				switch integration.Integration.Service {
				case SPOTIFY, VK, LASTFM:
					return integration.Enabled
				default:
					return false
				}
			},
		)

		lastFmIntegration, ok := lo.Find(
			integrations,
			func(integration model.ChannelsIntegrations) bool {
				return integration.Integration.Service == "LASTFM"
			},
		)

		var lfm *lastfm.LastFm
		if ok {
			lfm = lastfm.New(&lastFmIntegration)
		}

		spotifyIntegration, ok := lo.Find(
			integrations,
			func(integration model.ChannelsIntegrations) bool {
				return integration.Integration.Service == "SPOTIFY"
			},
		)
		var spoti *spotify.Spotify
		if ok {
			spoti = spotify.New(&spotifyIntegration, &db)
		}

		vkIntegration, ok := lo.Find(
			integrations,
			func(integration model.ChannelsIntegrations) bool {
				return integration.Integration.Service == "VK"
			},
		)
		var vk *vkIntegr.Vk
		if ok {
			vk = vkIntegr.New(&vkIntegration)
		}

		integrationsForFetch := lo.Map(
			integrations,
			func(integration model.ChannelsIntegrations, _ int) string {
				return integration.Integration.Service
			},
		)

		integrationsForFetch = append(integrationsForFetch, SOUNDTRACK)
		integrationsForFetch = append(integrationsForFetch, YOUTUBE_SR)

	checkServices:
		for _, integration := range integrationsForFetch {
			switch integration {
			case SPOTIFY:
				if spoti == nil {
					continue
				}
				track := spoti.GetTrack()
				if track != nil {
					result.Result = *track
					break checkServices
				}
			case LASTFM:
				if lfm == nil {
					continue
				}

				track := lfm.GetTrack()

				if track != nil {
					result.Result = *track
					break checkServices
				}
			case VK:
				if vk == nil {
					continue
				}
				track := vk.GetTrack()
				if track != nil {
					result.Result = *track
					break checkServices
				}
			case YOUTUBE_SR:
				redisData, err := redisClient.Get(
					context.Background(),
					fmt.Sprintf("songrequests:youtube:%s:currentPlaying", ctx.ChannelId),
				).Result()
				if err == redis.Nil {
					continue
				}
				if err != nil {
					zap.S().Error(err)
					continue
				}
				song := model.RequestedSong{}
				if err = db.Where("id = ?", redisData).First(&song).Error; err != nil {
					fmt.Println("song nog found", err)
					continue
				}

				result.Result = fmt.Sprintf(
					`"%s" youtu.be/%s requested by @%s`,
					song.Title,
					song.VideoID,
					song.OrderedByName,
				)
				break checkServices
			case SOUNDTRACK:
				cfg := do.MustInvoke[config.Config](di.Provider)
				tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

				twitchClient, err := twitch.NewAppClient(cfg, tokensGrpc)
				if err != nil {
					continue
				}
				tracks, err := twitchClient.GetSoundTrackCurrentTrack(&helix.SoundtrackCurrentTrackParams{
					BroadcasterID: ctx.ChannelId,
				})
				if err != nil {
					zap.S().Error(err)
					continue
				}

				if len(tracks.Data.Tracks) == 0 {
					continue
				}

				track := tracks.Data.Tracks[0]
				artists := lo.Map(track.Track.Artists, func(artist helix.SoundtrackTrackArtist, _ int) string {
					return artist.Name
				})
				result.Result = fmt.Sprintf("%s — %s", strings.Join(artists, ", "), track.Track.Title)
				break checkServices
			}
		}

		return result, nil
	},
}
