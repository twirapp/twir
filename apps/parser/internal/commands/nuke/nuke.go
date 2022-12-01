package nuke

import (
	"context"
	"fmt"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"

	"github.com/samber/lo"
)

var Command = types.DefaultCommand{
	Command: types.Command{
		Name: "nuke",
		Description: lo.ToPtr(
			"Mass remove messages in chat by message content. Usage: <b>!nuke phrase</b>",
		),
		Permission: "MODERATOR",
		Visible:    false,
		Module:     lo.ToPtr("CHANNEL"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		messages := []model.ChannelChatMessage{}

		if ctx.Text == nil {
			return nil
		}

		err := ctx.Services.Db.
			Where(
				`"canBeDeleted" = ? AND text LIKE ?`,
				true,
				"%"+strings.ToLower(*ctx.Text)+"%",
			).
			Find(&messages).
			Error
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if len(messages) == 0 {
			return nil
		}

		messages = lo.Filter(messages, func(m model.ChannelChatMessage, _ int) bool {
			return m.CanBeDeleted
		})
		mappedMessages := lo.Map(messages, func(m model.ChannelChatMessage, _ int) string {
			return m.MessageId
		})

		ctx.Services.BotsGrpc.DeleteMessage(context.Background(), &bots.DeleteMessagesRequest{
			ChannelId:   ctx.ChannelId,
			MessageIds:  mappedMessages,
			ChannelName: ctx.ChannelName,
		})

		ctx.Services.Db.Where(`"messageId" IN ?`, mappedMessages).
			Delete(&model.ChannelChatMessage{})

		return nil
	},
}
