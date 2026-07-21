package events

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/entities/platform"
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

	channelID, err := uuid.Parse(data.ChannelDBID)
	if err != nil {
		return fmt.Errorf("parse channel id: %w", err)
	}

	request := bots.SendMessageRequest{
		ChannelID:  channelID,
		Message:    msg,
		IsAnnounce: operation.UseAnnounce,
	}
	if data.Platform != "" {
		request.Platforms = []platform.Platform{data.Platform}
	}

	if err = c.bus.Bots.SendMessage.Publish(ctx, request); err != nil {
		return fmt.Errorf("cannot send message %w", err)
	}

	return nil
}
