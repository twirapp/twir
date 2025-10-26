package user

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/apps/parser/pkg/helpers"
	"github.com/twirapp/twir/libs/i18n"
)

var Age = &types.Variable{
	Name:         "user.age",
	Description:  lo.ToPtr("User account age"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		targetUserId := lo.
			IfF(
				len(parseCtx.Mentions) > 0, func() string {
					return parseCtx.Mentions[0].UserId
				},
			).
			Else(parseCtx.Sender.ID)
		user, err := parseCtx.Cacher.GetTwitchUserById(ctx, targetUserId)
		if err != nil {
			return nil, err
		}

		if user == nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.User.Errors.FindUserOnTwitch)
		} else {
			result.Result = helpers.Duration(
				user.CreatedAt.Time,
				&helpers.DurationOpts{
					UseUtc: true,
					Hide: helpers.DurationOptsHide{
						Seconds: true,
					},
				},
			)
		}

		return &result, nil
	},
}
