package brb

import (
	"context"
	"strconv"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/websockets"
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
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := types.CommandsHandlerResult{}

		if parseCtx.Text == nil {
			result.Result = []string{"you have to type minutes and optionally text"}
			return &result, nil
		}

		params := strings.Split(*parseCtx.Text, " ")

		if len(params) < 1 {
			result.Result = []string{"you have to type minutes and optionally text"}
			return &result, nil
		}

		minutes, err := strconv.Atoi(params[0])
		if err != nil {
			result.Result = []string{"first argument should be minutes"}
			return &result, nil
		}

		if minutes > 99999 || minutes <= 0 {
			result.Result = []string{"minutes cannot be more than 99999 and fewer then 1"}
			return &result, nil
		}

		text := strings.Join(params[1:], " ")
		if _, err := parseCtx.Services.GrpcClients.WebSockets.TriggerShowBrb(
			ctx,
			&websockets.TriggerShowBrbRequest{
				ChannelId: parseCtx.Channel.ID,
				Minutes:   int32(minutes),
				Text:      &text,
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
