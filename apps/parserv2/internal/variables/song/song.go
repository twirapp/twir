package song

import (
	lastfm "tsuwari/parser/internal/integrations/lastfm"
	model "tsuwari/parser/internal/models"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

const Name = "song"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := &types.VariableHandlerResult{}

	integrations := *ctx.GetEnabledIntegrations()
	if integrations == nil {
		result.Result = "you haven't enabled integrations for fetching song"
		return result, nil
	}

	integrations = lo.Filter(integrations, func(integration model.ChannelInegrationWithRelation, _ int) bool {
		switch integration.Integration.Service {
		case "LASTFM", "SPOTIFY", "VK":
			return integration.Enabled
		default:
			return false
		}
	})

	lastFmIntegration, _ := lo.Find(integrations, func(integration model.ChannelInegrationWithRelation) bool {
		return integration.Integration.Service == "LASTFM"
	})
	lfm := lastfm.New(&lastFmIntegration)

checkServices:
	for _, integration := range integrations {
		switch integration.Integration.Service {
		case "LASTFM":
			if lfm == nil {
				continue
			}

			track := lfm.GetRecentTrack()
			if track != nil {
				result.Result = *track
				break checkServices
			}
		}
	}

	return result, nil
}
