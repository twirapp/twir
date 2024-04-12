package messagehandler

import (
	"context"
	"fmt"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *MessageHandler) handleIncrementStreamMessages(
	ctx context.Context,
	msg handleMessage,
) error {
	fmt.Println("Incrementing stream messages", msg.DbStream)
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
