package userfollowage

import (
	types "tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "user.followage",
	Description: lo.ToPtr("User followage"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		follow := ctx.GetFollowAge()
		if follow == nil {
			result.Result = "not a follower"
		} else {
			result.Result = helpers.Duration(follow.FollowedAt)
		}

		return &result, nil
	},
}
