package random

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
)

var OnlineUser = &types.Variable{
	Name:                "random.online.user",
	Description:         lo.ToPtr("Choose random online user"),
	CanBeUsedInRegistry: true,
	NotCachable:         true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		randomUser, err := parseCtx.Services.UsersRepo.GetRandomOnlineUser(
			ctx,
			usersrepository.GetRandomOnlineUserInput{
				ChannelID: parseCtx.Channel.ID,
			},
		)
		if err != nil {
			return result, &types.CommandHandlerError{
				Message: "cannot get online user",
				Err:     err,
			}
		}

		result.Result = randomUser.UserName
		return result, nil
	},
}
