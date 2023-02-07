package user_follow

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	types "github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	"time"

	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
)

var FollowsinceVariable = types.Variable{
	Name:         "user.followsince",
	Description:  lo.ToPtr(`User follow since in "16 January 2023 (22 days)" format.`),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		cfg := do.MustInvoke[config.Config](di.Provider)
		tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

		twitchClient, err := twitch.NewAppClient(cfg, tokensGrpc)

		if err != nil {
			return nil, err
		}

		result := &types.VariableHandlerResult{}

		targetId := ctx.SenderId
		if ctx.Text != nil {
			users, err := twitchClient.GetUsers(&helix.UsersParams{
				Logins: []string{*ctx.Text},
			})

			if err != nil || len(users.Data.Users) == 0 {
				result.Result = "Cannot find user " + *ctx.Text + " on twitch."
				return result, nil
			}

			targetId = users.Data.Users[0].ID
		}

		if ctx.ChannelId == targetId {
			result.Result = "🎙️ broadcaster"
			return result, nil
		}

		follow := ctx.GetFollowAge(targetId)
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
