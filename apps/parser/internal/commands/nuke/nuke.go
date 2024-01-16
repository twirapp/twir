package nuke

import (
	"context"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/bots"

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
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		var messages []model.ChannelChatMessage

		if parseCtx.Text == nil {
			return nil, nil
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
			return nil, &types.CommandHandlerError{
				Message: "cannot get messages",
				Err:     err,
			}
		}

		if len(messages) == 0 {
			return nil, nil
		}

		messages = lo.Filter(
			messages, func(m model.ChannelChatMessage, _ int) bool {
				return m.CanBeDeleted
			},
		)
		mappedMessages := lo.Map(
			messages, func(m model.ChannelChatMessage, _ int) string {
				return m.MessageId
			},
		)

		if _, err = parseCtx.Services.GrpcClients.Bots.DeleteMessage(
			ctx,
			&bots.DeleteMessagesRequest{
				ChannelId:   parseCtx.Channel.ID,
				MessageIds:  mappedMessages,
				ChannelName: parseCtx.Channel.Name,
			},
		); err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot delete messages",
				Err:     err,
			}
		}

		if err = parseCtx.Services.Gorm.WithContext(ctx).Where(`"messageId" IN ?`, mappedMessages).
			Delete(&model.ChannelChatMessage{}).Error; err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot delete messages from db",
				Err:     err,
			}
		}

		return nil, nil
	},
}
