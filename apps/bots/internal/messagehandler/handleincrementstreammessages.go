package messagehandler

import (
	"context"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *MessageHandler) handleIncrementStreamMessages(
	ctx context.Context,
	msg handleMessage,
) error {
	if msg.DbStream == nil {
		return nil
	}

	if err := c.gorm.WithContext(ctx).Model(&model.ChannelsStreams{}).Update(
		"parsedMessages",
		msg.DbStream.ParsedMessages+1,
	).Error; err != nil {
		return err
	}

	return nil
}
