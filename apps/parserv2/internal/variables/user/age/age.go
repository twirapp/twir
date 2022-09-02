package userage

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"
)

const Name = "user.age"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	if ctx.Cache.TwitchUser == nil {
		result.Result = "error on getting user"
	} else {
		result.Result = helpers.Duration(ctx.Cache.TwitchUser.CreatedAt.Time)
	}

	return &result, nil
}
