package userage

import (
	"fmt"
	types "tsuwari/parser/internal/types"
	"tsuwari/parser/pkg/helpers"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/nicklaw5/helix"
	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "user.age",
	Description: lo.ToPtr("User account age"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		var user *helix.User
		if ctx.Text != nil {
			users, err := ctx.Services.Twitch.Client.GetUsers(&helix.UsersParams{
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
			result.Result = helpers.Duration(user.CreatedAt.Time)
		}

		return &result, nil
	},
}
