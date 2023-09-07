package user

import (
	"context"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/pkg/helpers"
	"github.com/satont/twir/libs/twitch"
)

var FollowAge = &types.Variable{
	Name:         "user.followage",
	Description:  lo.ToPtr(`User followage duration in "1y 3mo 22d" format`),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		twitchClient, err := twitch.NewAppClientWithContext(
			ctx,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			return nil, err
		}

		result := &types.VariableHandlerResult{}

		targetId := parseCtx.Sender.ID
		if parseCtx.Text != nil {
			userName := strings.ReplaceAll(*parseCtx.Text, "@", "")

			users, err := twitchClient.GetUsers(
				&helix.UsersParams{
					Logins: []string{userName},
				},
			)

			if err != nil || len(users.Data.Users) == 0 {
				result.Result = "Cannot find user " + *parseCtx.Text + " on twitch."
				return result, nil
			}

			targetId = users.Data.Users[0].ID
		}

		if parseCtx.Channel.ID == targetId {
			result.Result = "🎙️ broadcaster"
			return result, nil
		}

		follow := parseCtx.Cacher.GetTwitchUserFollow(ctx, targetId)
		if follow == nil {
			result.Result = "not a follower"
		} else {
			result.Result = helpers.Duration(
				follow.Followed.Time,
				&helpers.DurationOpts{
					UseUtc: true,
					Hide: helpers.DurationOptsHide{
						Minutes: true,
						Seconds: true,
					},
				},
			)
		}

		return result, nil
	},
}
