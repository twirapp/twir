package userage

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"

	"github.com/nicklaw5/helix"
)

const Name = "user.age"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	users, err := ctx.Services.Twitch.Client.GetUsers(&helix.UsersParams{
		IDs: []string{ctx.Context.SenderId},
	})

	if err != nil {
		return nil, err
	}

	if len(users.Data.Users) == 0 {
		result.Result = "not a follower"
	} else {
		result.Result = helpers.Duration(users.Data.Users[0].CreatedAt.Time)
	}

	return &result, nil
}
