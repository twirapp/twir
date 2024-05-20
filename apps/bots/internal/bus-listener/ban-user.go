package bus_listener

import (
	"context"
	"log/slog"

	"github.com/satont/twir/apps/bots/internal/twitchactions"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
)

func (c *BusListener) banUser(
	ctx context.Context,
	req bots.BanRequest,
) struct{} {
	if req.ChannelID == req.UserID {
		return struct{}{}
	}

	channelEntity := model.Channels{}
	if err := c.gorm.WithContext(ctx).Where(
		`"id" = ?`,
		req.ChannelID,
	).First(&channelEntity).Error; err != nil {
		c.logger.Error("cannot get channel entity", slog.Any("err", err))
		return struct{}{}
	}

	if err := c.twitchActions.Ban(
		ctx,
		twitchactions.BanOpts{
			Duration:       req.BanTime,
			Reason:         req.Reason,
			BroadcasterID:  req.ChannelID,
			UserID:         req.UserID,
			ModeratorID:    channelEntity.BotID,
			IsModerator:    req.IsModerator,
			AddModAfterBan: req.AddModAfterBan,
		},
	); err != nil {
		c.logger.Error("cannot ban user", slog.Any("err", err))
	}

	return struct{}{}
}
