package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/libs/twitch"
)

var FollowSince = &types.Variable{
	Name:         "user.followsince",
	Description:  lo.ToPtr(`User follow since in "16 January 2023 (22 days)" format.`),
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
			result.Result = fmt.Sprintf(
				"%s (%.0f days)",
				follow.Followed.Time.UTC().Format("2 January 2006"),
				time.Now().UTC().Sub(follow.Followed.Time.UTC()).Hours()/24,
			)
		}

		return result, nil
	},
}
