package nuke

import (
	"context"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"

	"go.uber.org/zap"
	"strings"

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
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		var messages []model.ChannelChatMessage

		if parseCtx.Text == nil {
			return nil
		}

		err := parseCtx.Services.Gorm.WithContext(ctx).
			Where(
				`"canBeDeleted" = ? AND text LIKE ?`,
				true,
				"%"+strings.ToLower(*parseCtx.Text)+"%",
			).
			Find(&messages).
			Error
		if err != nil {
			zap.S().Error(err)
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

		parseCtx.Services.GrpcClients.Bots.DeleteMessage(context.Background(), &bots.DeleteMessagesRequest{
			ChannelId:   parseCtx.Channel.ID,
			MessageIds:  mappedMessages,
			ChannelName: parseCtx.Channel.Name,
		})

		parseCtx.Services.Gorm.WithContext(ctx).Where(`"messageId" IN ?`, mappedMessages).
			Delete(&model.ChannelChatMessage{})

		return nil
	},
}
