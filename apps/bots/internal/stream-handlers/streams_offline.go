package stream_handlers

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/bus-core/twitch"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *PubSubHandlers) streamsOffline(
	ctx context.Context,
	data twitch.StreamOfflineMessage,
) (struct{}, error) {
	channel := model.Channels{}
	if err := c.db.WithContext(ctx).Where("id = ?", data.ChannelID).Find(&channel).Error; err != nil {
		c.logger.Error("cannot find channel", slog.String("channelId", data.ChannelID))
		return struct{}{}, err
	}

	if channel.ID == "" {
		return struct{}{}, nil
	}

	if err := c.greetingsCacher.Invalidate(ctx, channel.ID); err != nil {
		c.logger.Error(
			"cannot invalidate greetings cache",
			slog.String("channelId", data.ChannelID),
			slog.Any("err", err),
		)
		return struct{}{}, err
	}

	return struct{}{}, nil
}
