package user

import (
	"context"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/pkg/helpers"
)

var Age = &types.Variable{
	Name:         "user.age",
	Description:  lo.ToPtr("User account age"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		var user *helix.User
		if parseCtx.Text != nil {
			userName := strings.ReplaceAll(*parseCtx.Text, "@", "")

			cachedUser, err := parseCtx.Cacher.GetTwitchUserByName(ctx, userName)
			if err != nil {
				return nil, err
			}

			if cachedUser == nil {
				user = cachedUser
			}
		} else {
			user = parseCtx.Cacher.GetTwitchSenderUser(ctx)
		}

		if user == nil {
			result.Result = "Cannot find user on twitch."
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
