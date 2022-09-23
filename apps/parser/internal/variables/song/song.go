package song

import (
	model "tsuwari/models"
	lastfm "tsuwari/parser/internal/integrations/lastfm"
	spotify "tsuwari/parser/internal/integrations/spotify"
	vkIntegr "tsuwari/parser/internal/integrations/vk"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "currentsong",
	Description: lo.ToPtr("Current played song"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		integrations := *ctx.GetEnabledIntegrations()
		if integrations == nil {
			result.Result = "you haven't enabled integrations for fetching song"
			return result, nil
		}

		integrations = lo.Filter(integrations, func(integration model.ChannelInegrationWithRelation, _ int) bool {
			switch integration.Integration.Service {
			case "SPOTIFY", "LASTFM", "VK":
				return integration.Enabled
			default:
				return false
			}
		})

		lastFmIntegration, _ := lo.Find(integrations, func(integration model.ChannelInegrationWithRelation) bool {
			return integration.Integration.Service == "LASTFM"
		})
		lfm := lastfm.New(&lastFmIntegration)

		spotifyIntegration, _ := lo.Find(integrations, func(integration model.ChannelInegrationWithRelation) bool {
			return integration.Integration.Service == "SPOTIFY"
		})
		spoti := spotify.New(&spotifyIntegration, ctx.Services.Db)

		vkIntegration, _ := lo.Find(integrations, func(integration model.ChannelInegrationWithRelation) bool {
			return integration.Integration.Service == "VK"
		})
		vk := vkIntegr.New(&vkIntegration)

	checkServices:
		for _, integration := range integrations {
			switch integration.Integration.Service {
			case "SPOTIFY":
				if spoti == nil {
					continue
				}
				track := spoti.GetTrack()
				if track != nil {
					result.Result = *track
					break checkServices
				}
			case "LASTFM":
				if lfm == nil {
					continue
				}

				track := lfm.GetTrack()
				if track != nil {
					result.Result = *track
					break checkServices
				}
			case "VK":
				if vk == nil {
					continue
				}
				track := vk.GetTrack()
				if track != nil {
					result.Result = *track
					break checkServices
				}
			}
		}

		return result, nil
	},
}
