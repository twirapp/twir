package seventv

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var ProfileLink = &types.Variable{
	Name:         "7tv.profile.link",
	Description:  lo.ToPtr("Link to 7tv profile"),
	CommandsOnly: false,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile, err := parseCtx.Cacher.GetSeventvProfileGetTwitchId(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Seventv.Errors.ProfileNotFound.SetVars(locales.KeysVariablesSeventvErrorsProfileNotFoundVars{Reason: err.Error()}))
		} else {
			result.Result = "https://7tv.app/users/" + profile.Id
		}

		return &result, nil
	},
}
