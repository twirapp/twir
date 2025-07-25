package messagehandler

import (
	"context"
	"time"

	"github.com/twirapp/twir/libs/redis_keys"
	"github.com/twirapp/twir/libs/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *MessageHandler) handleIncrementStreamMessages(
	ctx context.Context,
	msg handleMessage,
) error {
	span := trace.SpanFromContext(ctx)
  defer span.End()
  span.SetAttributes(attribute.String("function.name", utils.GetFuncName()))

	if msg.EnrichedData.ChannelStream == nil {
		return nil
	}

	err := c.redis.Incr(
		ctx,
		redis_keys.StreamParsedMessages(
			msg.EnrichedData.ChannelStream.ID,
		),
	).Err()
	if err != nil {
		return err
	}

	return c.redis.Expire(
		ctx,
		redis_keys.StreamParsedMessages(msg.EnrichedData.ChannelStream.ID),
		50*time.Hour,
	).Err()
}
