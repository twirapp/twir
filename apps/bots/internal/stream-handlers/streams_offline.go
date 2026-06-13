package stream_handlers

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/logger"
)

func (c *PubSubHandlers) streamsOffline(
	ctx context.Context,
	data twitch.StreamOfflineMessage,
) (struct{}, error) {
	channel, found, err := c.findTwitchChannelByPlatformUserID(ctx, data.ChannelID)
	if err != nil {
		c.logger.Error("cannot find channel", slog.String("channelId", data.ChannelID))
		return struct{}{}, err
	}

	if !found {
		return struct{}{}, nil
	}

	if err := c.greetingsCacher.Invalidate(ctx, channel.ID); err != nil {
		c.logger.Error(
			"cannot invalidate greetings cache",
			slog.String("channelId", data.ChannelID),
			logger.Error(err),
		)
		return struct{}{}, err
	}

	return struct{}{}, nil
}
