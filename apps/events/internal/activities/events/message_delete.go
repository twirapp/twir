package events

import (
	"context"
	"fmt"

	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) MessageDelete(
	ctx context.Context,
	operation model.EventOperation,
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
