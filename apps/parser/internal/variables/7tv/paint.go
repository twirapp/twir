package seventv

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var Paint = &types.Variable{
	Name:         "7tv.profile.paint",
	Description:  lo.ToPtr("Paint of profile"),
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
		if profile.Style.ActivePaint == nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Seventv.Errors.NoPaint)
			return &result, nil
		}

		result.Result = profile.Style.ActivePaint.Name

		return &result, nil
	},
}

var UnlockedPaints = &types.Variable{
	Name:         "7tv.profile.unlockedpaints",
	Description:  lo.ToPtr("Num of unlocked paints"),
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

		result.Result = strconv.Itoa(len(profile.Inventory.Paints))

		return &result, nil
	},
}
