package messagehandler

import (
	"context"

	"github.com/twirapp/twir/libs/logger"
	chatmessages "github.com/twirapp/twir/libs/repositories/chat_messages"
	"github.com/twirapp/twir/libs/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *MessageHandler) handleSaveMessageBatched(ctx context.Context, data []handleMessage) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.String("function.name", utils.GetFuncName()))

	createMessageInputs := make([]chatmessages.CreateInput, len(data))

	for index, msg := range data {
		createMessageInputs[index] = chatmessages.CreateInput{
			ID:              msg.ID,
			ChannelID:       msg.BroadcasterUserId,
			UserID:          msg.ChatterUserId,
			Text:            msg.Message.Text,
			UserName:        msg.ChatterUserLogin,
			UserDisplayName: msg.ChatterUserName,
			UserColor:       msg.Color,
		}
	}

	err := c.chatMessagesRepository.CreateMany(ctx, createMessageInputs)
	if err != nil {
		c.logger.Error("cannot save chat messages to db", logger.Error(err))
	}
}

func (c *MessageHandler) handleSaveMessage(
	_ context.Context,
	msg handleMessage,
) error {
	if msg.Message == nil {
		return nil
	}

	c.messagesSaveBatcher.Add(msg)

	return nil
}
