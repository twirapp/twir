package user

import (
	"context"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	"github.com/satont/tsuwari/apps/parser/pkg/helpers"
	"github.com/satont/tsuwari/libs/twitch"
)

var Age = &types.Variable{
	Name:         "user.age",
	Description:  lo.ToPtr("User account age"),
	CommandsOnly: true,
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		twitchClient, err := twitch.NewAppClientWithContext(
			ctx,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			return nil, err
		}

		result := types.VariableHandlerResult{}

		var user *helix.User
		if parseCtx.Text != nil {
			userName := strings.ReplaceAll(*parseCtx.Text, "@", "")

			users, err := twitchClient.GetUsers(&helix.UsersParams{
				Logins: []string{userName},
			})

			if err == nil && len(users.Data.Users) != 0 {
				user = &users.Data.Users[0]
			}
		} else {
			user = parseCtx.Cacher.GetTwitchSenderUser(ctx)
		}

		if user == nil {
			result.Result = "Cannot find user on twitch."
		} else {
			result.Result = helpers.Duration(user.CreatedAt.Time, &helpers.DurationOpts{
				UseUtc: true,
				Hide: helpers.DurationOptsHide{
					Seconds: true,
				},
			})
		}

		return &result, nil
	},
}
