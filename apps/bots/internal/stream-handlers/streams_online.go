package stream_handlers

import (
	"context"
	"log/slog"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/repositories/greetings"
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

	err := c.greetingsRepository.UpdateManyByChannelID(
		ctx, greetings.UpdateManyInput{
			ChannelID: channel.ID,
			Processed: lo.ToPtr(false),
		},
	)

	if err != nil {
		c.logger.Error(
			"cannot update channel greetings",
			slog.String("channelId", data.ChannelID),
			slog.Any("err", err),
		)
	}

	return struct{}{}
}
