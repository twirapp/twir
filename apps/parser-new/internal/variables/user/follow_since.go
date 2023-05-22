package user

import (
	"context"
	"fmt"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
	"github.com/satont/tsuwari/libs/twitch"
)

var FollowSince = &types.Variable{
	Name:         "user.followsince",
	Description:  lo.ToPtr(`User follow since in "16 January 2023 (22 days)" format.`),
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

		result := &types.VariableHandlerResult{}

		targetId := parseCtx.Sender.ID
		if parseCtx.Text != nil {
			users, err := twitchClient.GetUsers(&helix.UsersParams{
				Logins: []string{*parseCtx.Text},
			})

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
			result.Result = fmt.Sprintf(
				"%s (%.0f days)",
				follow.FollowedAt.UTC().Format("2 January 2006"),
				time.Now().UTC().Sub(follow.FollowedAt.UTC()).Hours()/24,
			)
		}

		return result, nil
	},
}
