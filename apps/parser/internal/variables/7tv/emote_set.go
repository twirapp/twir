package seventv

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var EmoteSetName = &types.Variable{
	Name:         "7tv.emoteset.name",
	Description:  lo.ToPtr("Name of 7tv emote set"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile, err := parseCtx.Services.SevenTvCache.Get(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = fmt.Sprintf("[Twir err] Failed to get 7tv profile: %s", err)
			return &result, nil
		}
		if profile.EmoteSet == nil {
			result.Result = fmt.Sprintf(
				"[Twir err] Failed to get 7tv emote set: %s",
				"emote set is not set",
			)
			return &result, nil
		}

		result.Result = profile.EmoteSet.Name

		return &result, nil
	},
}

var EmoteSetLink = &types.Variable{
	Name:         "7tv.emoteset.link",
	Description:  lo.ToPtr("Link to 7tv emote set"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile, err := parseCtx.Services.SevenTvCache.Get(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = fmt.Sprintf("[Twir err] Failed to get 7tv profile: %s", err)
			return &result, nil
		}
		if profile.EmoteSet == nil {
			result.Result = fmt.Sprintf(
				"[Twir err] Failed to get 7tv emote set: %s",
				"emote set is not set",
			)
			return &result, nil
		}

		result.Result = fmt.Sprintf("https://7tv.app/emote-sets/%s", profile.EmoteSet.Id)

		return &result, nil
	},
}

var EmoteSetCount = &types.Variable{
	Name:         "7tv.emoteset.emotes.count",
	Description:  lo.ToPtr("Count of emotes in emote set"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile, err := parseCtx.Services.SevenTvCache.Get(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = fmt.Sprintf("[Twir err] Failed to get 7tv profile: %s", err)
			return &result, nil
		}
		if profile.EmoteSet == nil {
			result.Result = fmt.Sprintf(
				"[Twir err] Failed to get 7tv emote set: %s",
				"emote set is not set",
			)
			return &result, nil
		}

		result.Result = fmt.Sprint(len(profile.EmoteSet.Emotes))

		return &result, nil
	},
}
