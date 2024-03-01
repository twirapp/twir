package nuke

import (
	"context"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/libs/bus-core/bots"

	model "github.com/satont/twir/libs/gomodels"
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
		if parseCtx.Text == nil {
			return nil, nil
		}

		var messages []model.ChannelChatMessage
		err := parseCtx.Services.Gorm.WithContext(ctx).
			Where(
				`"canBeDeleted" IS TRUE AND text LIKE ? AND "createdAt" > NOW() - INTERVAL '60 minutes' AND "channelId" = ?`,
				"%"+strings.ToLower(*parseCtx.Text)+"%",
				parseCtx.Channel.ID,
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

		mappedMessagesIds := make([]string, 0, len(messages))

		for _, message := range messages {
			if !message.CanBeDeleted {
				continue
			}

			mappedMessagesIds = append(mappedMessagesIds, message.MessageId)
		}

		if err := parseCtx.Services.Bus.Bots.DeleteMessage.Publish(
			bots.DeleteMessageRequest{
				ChannelId:   parseCtx.Channel.ID,
				MessageIds:  mappedMessagesIds,
				ChannelName: &parseCtx.Channel.Name,
			},
		); err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot delete messages",
				Err:     err,
			}
		}

		if err = parseCtx.Services.Gorm.WithContext(ctx).Where(`"messageId" IN ?`, mappedMessagesIds).
			Delete(&model.ChannelChatMessage{}).Error; err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot delete messages from db",
				Err:     err,
			}
		}

		return nil, nil
	},
}
