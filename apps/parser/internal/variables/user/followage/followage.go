package userfollowage

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/libs/twitch"

	"github.com/satont/tsuwari/apps/parser/internal/di"
	types "github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"github.com/satont/tsuwari/apps/parser/pkg/helpers"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"

	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
)

var Variable = types.Variable{
	Name:         "user.followage",
	Description:  lo.ToPtr("User followage"),
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
			result.Result = helpers.Duration(follow.FollowedAt, &helpers.DurationOpts{
				UseUtc: true,
				Hide: helpers.DurationOptsHide{
					Minutes: true,
					Seconds: true,
				},
			})
		}

		return result, nil
	},
}
