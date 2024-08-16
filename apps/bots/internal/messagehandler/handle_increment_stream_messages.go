package messagehandler

import (
	"context"
	"time"

	"github.com/twirapp/twir/libs/redis_keys"
)

func (c *MessageHandler) handleIncrementStreamMessages(
	ctx context.Context,
	msg handleMessage,
) error {
	if msg.DbStream == nil {
		return nil
	}

	err := c.redis.Incr(
		ctx,
		redis_keys.StreamParsedMessages(
			msg.DbStream.ID,
		),
	).Err()
	if err != nil {
		return err
	}

	return c.redis.Expire(ctx, redis_keys.StreamParsedMessages(msg.DbStream.ID), 50*time.Hour).Err()
}
