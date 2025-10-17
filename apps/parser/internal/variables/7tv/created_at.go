package seventv

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/apps/parser/pkg/helpers"
	"github.com/twirapp/twir/libs/i18n"
)

var ProfileCreatedAt = &types.Variable{
	Name:         "7tv.profile.createdAt",
	Description:  lo.ToPtr("Date when profile created on 7tv"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile, err := parseCtx.Cacher.GetSeventvProfileGetTwitchId(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = i18n.GetCtx(
				ctx,
				locales.Translations.Variables.Seventv.Errors.ProfileNotFound.SetVars(locales.KeysVariablesSeventvErrorsProfileNotFoundVars{Reason: err.Error()}),
			)
		} else {
			result.Result = helpers.Duration(profile.MainConnection.LinkedAt, &helpers.DurationOpts{})
		}

		return &result, nil
	},
}
