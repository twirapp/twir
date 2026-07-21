package channel

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var ID = &types.Variable{
	Name:         "channel.id",
	Description:  new("Internal ID of channel"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		return &types.VariableHandlerResult{Result: parseCtx.Channel.DBChannelID}, nil
	},
}

var TwitchID = &types.Variable{
	Name:         "channel.twitch.id",
	Description:  new("Twitch ID of channel. Empty in case of twitch not connected"),
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

		return &types.VariableHandlerResult{Result: *ch.TwitchPlatformID}, nil
	},
}

var KickID = &types.Variable{
	Name:         "channel.kick.id",
	Description:  new("Kick ID of channel. Empty in case of twitch not connected"),
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

		return &types.VariableHandlerResult{Result: *ch.KickPlatformID}, nil
	},
}
