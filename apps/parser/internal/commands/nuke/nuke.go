package nuke

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"gorm.io/gorm"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"

	"github.com/samber/lo"
)

var Command = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name: "nuke",
		Description: null.StringFrom(
			"Mass remove messages in chat by message content. Usage: <b>!nuke phrase</b>",
		),
		RolesIDS: pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:   "MODERATION",
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		botsGrpc := do.MustInvoke[bots.BotsClient](di.Provider)
		db := do.MustInvoke[gorm.DB](di.Provider)

		messages := []model.ChannelChatMessage{}

		if ctx.Text == nil {
			return nil
		}

		err := db.
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

		botsGrpc.DeleteMessage(context.Background(), &bots.DeleteMessagesRequest{
			ChannelId:   ctx.ChannelId,
			MessageIds:  mappedMessages,
			ChannelName: ctx.ChannelName,
		})

		db.Where(`"messageId" IN ?`, mappedMessages).
			Delete(&model.ChannelChatMessage{})

		return nil
	},
}
