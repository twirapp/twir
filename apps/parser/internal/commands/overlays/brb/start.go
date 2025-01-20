package brb

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
)

const (
	startTimeArgName = "time"
	startTextArgName = "text"
)

var Start = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "brb",
		Description: null.StringFrom("Be right back overlay start command"),
		Module:      "OVERLAYS",
		IsReply:     true,
		Visible:     false,
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeBroadcaster.String(),
		},
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name: startTimeArgName,
			Min:  lo.ToPtr(1),
			Max:  lo.ToPtr(99999),
		},
		command_arguments.VariadicString{
			Name:     startTextArgName,
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := types.CommandsHandlerResult{}

		textArg := parseCtx.ArgsParser.Get(startTextArgName)
		var text *string
		if textArg != nil {
			text = lo.ToPtr(textArg.String())
		}

		if _, err := parseCtx.Services.GrpcClients.WebSockets.TriggerShowBrb(
			ctx,
			&websockets.TriggerShowBrbRequest{
				ChannelId: parseCtx.Channel.ID,
				Minutes:   int32(parseCtx.ArgsParser.Get(startTimeArgName).Int()),
				Text:      text,
			},
		); err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot trigger show brb",
				Err:     err,
			}
		}

		return &result, nil
	},
}
