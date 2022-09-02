package userfollowage

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"
)

const Name = "user.followage"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	if ctx.Cache.TwitchFollow == nil {
		result.Result = "not a follower"
	} else {
		result.Result = helpers.Duration(ctx.Cache.TwitchFollow.FollowedAt)
	}

	return &result, nil
}
