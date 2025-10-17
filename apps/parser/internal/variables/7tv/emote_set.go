package seventv

import (
	"context"
	"fmt"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var EmoteSetName = &types.Variable{
	Name:         "7tv.emoteset.name",
	Description:  lo.ToPtr("Name of 7tv emote set"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile, err := parseCtx.Cacher.GetSeventvProfileGetTwitchId(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Seventv.Errors.ProfileNotFound.SetVars(locales.KeysVariablesSeventvErrorsProfileNotFoundVars{Reason: err.Error()}))
			return &result, nil
		}
		if profile.Style.ActiveEmoteSet == nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Seventv.Errors.NoActiveSet)
			return &result, nil
		}

		result.Result = profile.Style.ActiveEmoteSet.Name

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

		profile, err := parseCtx.Cacher.GetSeventvProfileGetTwitchId(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Seventv.Errors.ProfileNotFound.SetVars(locales.KeysVariablesSeventvErrorsProfileNotFoundVars{Reason: err.Error()}))
			return &result, nil
		}
		if profile.Style.ActiveEmoteSet == nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Seventv.Errors.NoActiveSet)
			return &result, nil
		}

		result.Result = fmt.Sprintf("https://7tv.app/emote-sets/%s", profile.Style.ActiveEmoteSet.Id)

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

		profile, err := parseCtx.Cacher.GetSeventvProfileGetTwitchId(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Seventv.Errors.ProfileNotFound.SetVars(locales.KeysVariablesSeventvErrorsProfileNotFoundVars{Reason: err.Error()}))
			return &result, nil
		}
		if profile.Style.ActiveEmoteSet == nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Seventv.Errors.NoActiveSet)
			return &result, nil
		}

		result.Result = fmt.Sprint(len(profile.Style.ActiveEmoteSet.Emotes.Items))

		return &result, nil
	},
}

var EmoteSetCapacity = &types.Variable{
	Name:         "7tv.emoteset.capacity",
	Description:  lo.ToPtr("Capacity of set"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile, err := parseCtx.Cacher.GetSeventvProfileGetTwitchId(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Seventv.Errors.ProfileNotFound.SetVars(locales.KeysVariablesSeventvErrorsProfileNotFoundVars{Reason: err.Error()}))
			return &result, nil
		}
		if profile.Style.ActiveEmoteSet == nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Seventv.Errors.NoActiveSet)
			return &result, nil
		}
		if profile.Style.ActiveEmoteSet.Capacity == nil {
			result.Result = "0"
			return &result, nil
		}

		result.Result = strconv.Itoa(*profile.Style.ActiveEmoteSet.Capacity)

		return &result, nil
	},
}
