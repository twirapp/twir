package song

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	spotify "github.com/satont/tsuwari/libs/integrations/spotify"
)

var Song = &types.Variable{
	Name:        "currentsong",
	Description: lo.ToPtr("Print current played song from Spotify, Last.fm, e.t.c, and also from song requests."),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		integrations := parseCtx.Cacher.GetEnabledChannelIntegrations(ctx)

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
			spoti = spotify.New(spotifyIntegration, parseCtx.Services.Gorm)
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
					result.Result = *track
					break checkServices
				}
			case model.IntegrationServiceLastfm:
				if lfm == nil {
					continue
				}

				track := lfm.GetTrack()

				if track != nil {
					result.Result = *track
					break checkServices
				}
			case model.IntegrationServiceVK:
				if vk == nil {
					continue
				}
				track := vk.GetTrack(ctx)
				if track != nil {
					result.Result = *track
					break checkServices
				}
			case "YOUTUBE_SR":
				redisData, err := parseCtx.Services.Redis.Get(
					context.Background(),
					fmt.Sprintf("songrequests:youtube:%s:currentPlaying", parseCtx.Channel.ID),
				).Result()
				if err == redis.Nil {
					continue
				}
				if err != nil {
					parseCtx.Services.Logger.Sugar().Error(err)
					continue
				}
				song := model.RequestedSong{}
				if err = parseCtx.Services.Gorm.
					WithContext(ctx).
					Where("id = ?", redisData).
					First(&song).Error; err != nil {
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
			}
		}

		return result, nil
	},
}
