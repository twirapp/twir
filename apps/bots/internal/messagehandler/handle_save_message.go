package messagehandler

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/chat_messages"
)

func (c *MessageHandler) handleSaveMessage(
	ctx context.Context,
	msg handleMessage,
) error {
	if msg.Message == nil {
		return nil
	}

	err := c.chatMessagesRepository.Create(
		ctx,
		chat_messages.CreateInput{
			ChannelID:       msg.BroadcasterUserId,
			UserID:          msg.ChatterUserId,
			Text:            msg.Message.Text,
			UserName:        msg.ChatterUserLogin,
			UserDisplayName: msg.ChatterUserName,
			UserColor:       msg.Color,
		},
	)

	return err
}
