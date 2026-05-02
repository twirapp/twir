package stream_handlers

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/greetings"
)

func (c *PubSubHandlers) streamsOnline(
	ctx context.Context,
	data twitch.StreamOnlineMessage,
) (struct{}, error) {
	channel, found, err := c.findTwitchChannelByPlatformUserID(ctx, data.ChannelID)
	if err != nil {
		c.logger.Error(
			"cannot find channel",
			slog.String("channelId", data.ChannelID),
		)
		return struct{}{}, err
	}

	if !found {
		return struct{}{}, nil
	}

	channelID, err := uuid.Parse(channel.ID)
	if err != nil {
		return struct{}{}, err
	}

	err = c.greetingsRepository.UpdateManyByChannelID(
		ctx, greetings.UpdateManyInput{
			ChannelID: channelID,
			Processed: lo.ToPtr(false),
		},
	)
	if err != nil {
		c.logger.Error(
			"cannot update channel greetings",
			slog.String("channelId", data.ChannelID),
			logger.Error(err),
		)
		return struct{}{}, err
	}

	if err = c.greetingsCacher.Invalidate(ctx, channel.ID); err != nil {
		c.logger.Error(
			"cannot invalidate greetings cache",
			slog.String("channelId", data.ChannelID),
			logger.Error(err),
		)
		return struct{}{}, err
	}

	return struct{}{}, nil
}
