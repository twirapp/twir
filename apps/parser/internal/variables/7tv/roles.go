package seventv

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var Roles = &types.Variable{
	Name:         "7tv.roles",
	Description:  lo.ToPtr("Roles of user on 7tv"),
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

		roles := make([]string, 0, len(profile.Roles))
		for _, role := range profile.Roles {
			if role.Name == "Default" {
				continue
			}

			roles = append(roles, role.Name)
		}

		if len(roles) == 0 {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Seventv.Errors.NoRoles)
		} else {
			result.Result = strings.Join(roles, ", ")
		}

		return &result, nil
	},
}
