package userage

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"

	types "github.com/satont/tsuwari/apps/parser/internal/types"
	"github.com/satont/tsuwari/apps/parser/pkg/helpers"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
)

var Variable = types.Variable{
	Name:         "user.age",
	Description:  lo.ToPtr("User account age"),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		cfg := do.MustInvoke[config.Config](di.Provider)
		tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

		twitchClient, err := twitch.NewAppClient(cfg, tokensGrpc)

		if err != nil {
			return nil, err
		}

		result := types.VariableHandlerResult{}

		var user *helix.User
		if ctx.Text != nil {
			users, err := twitchClient.GetUsers(&helix.UsersParams{
				Logins: []string{*ctx.Text},
			})

			if err == nil && len(users.Data.Users) != 0 {
				user = &users.Data.Users[0]
			}
		} else {
			user = ctx.GetTwitchUser()
		}

		if user == nil {
			name := lo.If(ctx.Text != nil, *ctx.Text).Else(ctx.SenderName)
			result.Result = fmt.Sprintf("Cannot find user %s on twitch.", name)
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
