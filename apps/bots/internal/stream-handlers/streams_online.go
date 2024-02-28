package stream_handlers

import (
	"context"
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/twitch"
)

func (c *PubSubHandlers) streamsOnline(
	ctx context.Context,
	data twitch.StreamOnlineMessage,
) struct{} {
	channel := model.Channels{}
	if err := c.db.WithContext(ctx).Where("id = ?", data.ChannelID).Find(&channel).Error; err != nil {
		c.logger.Error(
			"cannot find channel",
			slog.String("channelId", data.ChannelID),
		)
		return struct{}{}
	}

	if channel.ID == "" {
		return struct{}{}
	}

	err := c.db.Model(&model.ChannelsGreetings{}).
		Where(`"channelId" = ?`, channel.ID).
		Update("processed", false).Error
	if err != nil {
		c.logger.Error(
			"cannot update channel greetings",
			slog.String("channelId", data.ChannelID),
			slog.Any("err", err),
		)
	}

	return struct{}{}
}
