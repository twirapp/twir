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

	input := chat_messages.CreateInput{
		ChannelID:       msg.BroadcasterUserId,
		UserID:          msg.ChatterUserId,
		Text:            msg.Message.Text,
		UserName:        msg.ChatterUserLogin,
		UserDisplayName: msg.ChatterUserName,
		UserColor:       msg.Color,
	}

	err := c.chatMessagesRepository.Create(
		ctx,
		input,
	)
	if err != nil {
		return err
	}

	return nil

	// handleSaveMessagesQueueLock.Lock()
	// defer handleSaveMessagesQueueLock.Unlock()
	//
	// handleSaveMessagesQueue = append(handleSaveMessagesQueue, msg)
	//
	// if len(handleSaveMessagesQueue) < 5 {
	// 	return nil
	// }
	//
	// inputs := make([]chat_messages.CreateInput, 0, len(handleSaveMessagesQueue))
	// for _, m := range handleSaveMessagesQueue {
	// 	inputs = append(
	// 		inputs,
	// 		chat_messages.CreateInput{
	// 			ChannelID:       m.BroadcasterUserId,
	// 			UserID:          m.ChatterUserId,
	// 			Text:            m.Message.Text,
	// 			UserName:        m.ChatterUserLogin,
	// 			UserDisplayName: m.ChatterUserName,
	// 			UserColor:       m.Color,
	// 		},
	// 	)
	// }
	//
	// err := c.chatMessagesRepository.CreateMany(
	// 	ctx,
	// 	inputs,
	// )
	// if err != nil {
	// 	return err
	// }
	//
	// handleSaveMessagesQueue = nil
	// return nil
}
