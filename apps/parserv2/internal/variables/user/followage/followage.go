package userfollowage

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"

	"github.com/nicklaw5/helix"
)

const Name = "user.followage"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	follow, err := ctx.Services.Twitch.Client.GetUsersFollows(&helix.UsersFollowsParams{
		FromID: ctx.Context.SenderId,
		ToID:   ctx.Context.ChannelId,
	})

	if err != nil {
		return nil, err
	}

	if len(follow.Data.Follows) == 0 {
		result.Result = "not a follower"
	} else {

		result.Result = helpers.Duration(follow.Data.Follows[0].FollowedAt)
	}

	return &result, nil
}
