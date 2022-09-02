package userfollowage

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"
)

const Name = "user.followage"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	follow := ctx.GetFollowAge()
	if follow == nil {
		result.Result = "not a follower"
	} else {
		result.Result = helpers.Duration(follow.FollowedAt)
	}

	return &result, nil
}
