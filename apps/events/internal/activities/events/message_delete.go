package events

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) MessageDelete(
	ctx context.Context,
	_ model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if data.ChatMessageId == "" {
		return fmt.Errorf("message id is empty: %s", data.ChatMessageId)
	}

	if err := c.bus.Bots.DeleteMessage.Publish(
		ctx,
		bots.DeleteMessageRequest{
			ChannelId:   data.ChannelID,
			ChannelName: nil,
			MessageIds:  []string{data.ChatMessageId},
		},
	); err != nil {
		return fmt.Errorf("cannot delete message %w", err)
	}

	return nil
}
