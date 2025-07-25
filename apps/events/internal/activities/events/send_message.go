package events

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/events/internal/shared"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) SendMessage(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	msg, err := c.hydrator.HydrateStringWithData(data.ChannelID, operation.Input.String, data)
	if err != nil {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	if err = c.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelId:  data.ChannelID,
			Message:    msg,
			IsAnnounce: operation.UseAnnounce,
		},
	); err != nil {
		return fmt.Errorf("cannot send message %w", err)
	}

	return nil
}
