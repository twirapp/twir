package seventv

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var EditorForCount = &types.Variable{
	Name:         "7tv.editorfor.count",
	Description:  lo.ToPtr("Count of channels where user is editor"),
	CommandsOnly: false,
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
			return &result, nil
		}

		result.Result = strconv.Itoa(len(profile.EditorFor))

		return &result, nil
	},
}
