package channel

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var TwitchName = &types.Variable{
	Name:         "channel.twitch.name",
	Description:  new("Twitch Name of channel. Empty in case of twitch not connected"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		uid, err := uuid.Parse(parseCtx.Channel.DBChannelID)
		if err != nil {
			return nil, err
		}

		ch, err := parseCtx.Services.ChannelService.GetChannelByID(ctx, uid)
		if err != nil {
			return nil, err
		}

		if !ch.TwitchConnected() {
			return &types.VariableHandlerResult{Result: ""}, nil
		}

		user, err := parseCtx.Services.UsersRepo.GetByID(ctx, *ch.TwitchUserID)
		if err != nil {
			return nil, err
		}

		if user.IsNil() {
			return &types.VariableHandlerResult{Result: ""}, nil
		}

		return &types.VariableHandlerResult{Result: user.Login}, nil
	},
}

var KickName = &types.Variable{
	Name:         "channel.kick.name",
	Description:  new("Kick Name of channel. Empty in case of Kick not connected"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		uid, err := uuid.Parse(parseCtx.Channel.DBChannelID)
		if err != nil {
			return nil, err
		}

		ch, err := parseCtx.Services.ChannelService.GetChannelByID(ctx, uid)
		if err != nil {
			return nil, err
		}

		if !ch.KickConnected() {
			return &types.VariableHandlerResult{Result: ""}, nil
		}

		user, err := parseCtx.Services.UsersRepo.GetByID(ctx, *ch.KickUserID)
		if err != nil {
			return nil, err
		}

		if user.IsNil() {
			return &types.VariableHandlerResult{Result: ""}, nil
		}

		return &types.VariableHandlerResult{Result: user.Login}, nil
	},
}
