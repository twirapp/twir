package messagehandler

import (
	"context"
	"sync"

	"github.com/twirapp/twir/libs/repositories/chat_messages"
)

var handleSaveMessagesQueue []handleMessage
var handleSaveMessagesQueueLock sync.Mutex

func (c *MessageHandler) handleSaveMessage(
	ctx context.Context,
	msg handleMessage,
) error {
	if msg.Message == nil {
		return nil
	}

	handleSaveMessagesQueueLock.Lock()
	defer handleSaveMessagesQueueLock.Unlock()

	handleSaveMessagesQueue = append(handleSaveMessagesQueue, msg)

	bufferSize := 5
	if c.config.AppEnv != "production" {
		bufferSize = 1
	}

	if len(handleSaveMessagesQueue) < bufferSize {
		return nil
	}

	inputs := make([]chat_messages.CreateInput, 0, len(handleSaveMessagesQueue))
	for _, m := range handleSaveMessagesQueue {
		inputs = append(
			inputs,
			chat_messages.CreateInput{
				ID:              m.ID,
				ChannelID:       m.BroadcasterUserId,
				UserID:          m.ChatterUserId,
				Text:            m.Message.Text,
				UserName:        m.ChatterUserLogin,
				UserDisplayName: m.ChatterUserName,
				UserColor:       m.Color,
			},
		)
	}

	err := c.chatMessagesRepository.CreateMany(
		ctx,
		inputs,
	)
	if err != nil {
		return err
	}

	handleSaveMessagesQueue = nil
	return nil
}
