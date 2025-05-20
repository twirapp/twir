package messagehandler

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/repositories/chat_messages"
)

// var handleSaveMessagesQueue []handleMessage
// var handleSaveMessagesQueueLock sync.Mutex

func (c *MessageHandler) handleSaveMessageBatched(ctx context.Context, data []handleMessage) {
	for _, msg := range data {
		err := c.chatMessagesRepository.Create(
			ctx,
			chat_messages.CreateInput{
				ID:              msg.ID,
				ChannelID:       msg.BroadcasterUserId,
				UserID:          msg.ChatterUserId,
				Text:            msg.Message.Text,
				UserName:        msg.ChatterUserLogin,
				UserDisplayName: msg.ChatterUserName,
				UserColor:       msg.Color,
			},
		)
		if err != nil {
			c.logger.Error("cannot save chat message to db", slog.Any("err", err))
			continue
		}
	}
}

func (c *MessageHandler) handleSaveMessage(
	ctx context.Context,
	msg handleMessage,
) error {
	if msg.Message == nil {
		return nil
	}

	c.messagesSaveBatcher.Add(msg)

	return nil
}
