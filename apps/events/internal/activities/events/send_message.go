package events

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) SendMessage(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Input == nil || *operation.Input == "" {
		return fmt.Errorf("input is required for send message operation")
	}

	msg, err := c.hydrator.HydrateStringWithData(data.ChannelID, data.ChannelTwitchUserID, data.ChannelDBID, *operation.Input, data)
	if err != nil {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	var internalChannelID *uuid.UUID
	if parsedChannelID, err := uuid.Parse(data.ChannelDBID); err == nil {
		internalChannelID = &parsedChannelID
	}

	if err = c.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelId:         data.ChannelID,
			InternalChannelID: internalChannelID,
			PlatformChannelID: data.ChannelID,
			Platform:          data.Platform.String(),
			Message:           msg,
			IsAnnounce:        operation.UseAnnounce,
		},
	); err != nil {
		return fmt.Errorf("cannot send message %w", err)
	}

	return nil
}
