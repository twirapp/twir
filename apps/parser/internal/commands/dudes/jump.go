package dudes

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
)

var Jump = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "jump",
		Description: null.StringFrom("Triggers jump of dude in dudes overlay"),
		RolesIDS:    pq.StringArray{},
		Module:      "DUDES",
		Visible:     true,
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		_, err := parseCtx.Services.GrpcClients.WebSockets.DudesJump(
			ctx, &websockets.DudesJumpRequest{
				ChannelId:       parseCtx.Channel.ID,
				UserId:          parseCtx.Sender.ID,
				UserDisplayName: parseCtx.Sender.DisplayName,
				UserName:        parseCtx.Sender.Name,
				UserColor:       parseCtx.Sender.Color,
			},
		)

		result := &types.CommandsHandlerResult{
			Result: []string{},
		}

		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = append(result.Result, "[Twir error] cannot trigger dudes jump")
		}

		return result, nil
	},
}
