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

	return c.gorm.
		WithContext(ctx).
		Model(&model.ChannelsStreams{}).
		Where(`"userId" = ?`, msg.BroadcasterUserId).
		Update(
			"parsedMessages",
			msg.DbStream.ParsedMessages+1,
		).Error
}
