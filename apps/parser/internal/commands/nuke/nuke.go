package nuke

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/libs/bus-core/bots"
	chat_messages_store "github.com/twirapp/twir/libs/bus-core/chat-messages-store"

	model "github.com/satont/twir/libs/gomodels"
)

const (
	nukePhraseArgName = "phrase"
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
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.VariadicString{
			Name: nukePhraseArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		phrase := parseCtx.ArgsParser.Get(nukePhraseArgName).String()

		messages, err := parseCtx.Services.Bus.ChatMessagesStore.GetChatMessagesByTextForDelete.Request(
			ctx,
			chat_messages_store.GetChatMessagesByTextRequest{
				ChannelID: parseCtx.Channel.ID,
				Text:      phrase,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get messages",
				Err:     err,
			}
		}

		if len(messages.Data.Messages) == 0 {
			return nil, nil
		}

		mappedMessagesIds := make([]string, 0, len(messages.Data.Messages))

		for _, message := range messages.Data.Messages {
			mappedMessagesIds = append(mappedMessagesIds, message.MessageID)
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

		redisMessagesIds := make([]string, 0, len(messages.Data.Messages))
		for _, message := range messages.Data.Messages {
			redisMessagesIds = append(redisMessagesIds, message.RedisID)
		}

		if err = parseCtx.Services.Bus.ChatMessagesStore.RemoveMessages.Publish(
			chat_messages_store.RemoveMessagesRequest{
				MessagesRedisIDS: redisMessagesIds,
			},
		); err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot remove messages",
				Err:     err,
			}
		}

		return nil, nil
	},
}
