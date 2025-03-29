package bus_listener

import (
	"context"
	"log/slog"
	"time"

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

func (c *BusListener) banUsers(
	ctx context.Context,
	req []bots.BanRequest,
) struct{} {
	uniqueChannelsIdsMap := make(map[string]struct{}, len(req))
	for _, r := range req {
		if r.ChannelID == r.UserID {
			continue
		}

		if _, ok := uniqueChannelsIdsMap[r.ChannelID]; !ok {
			uniqueChannelsIdsMap[r.ChannelID] = struct{}{}
		}
	}

	if len(uniqueChannelsIdsMap) == 0 {
		return struct{}{}
	}

	uniqueChannelsIds := make([]string, 0, len(uniqueChannelsIdsMap))
	for k := range uniqueChannelsIdsMap {
		uniqueChannelsIds = append(uniqueChannelsIds, k)
	}

	var channelsEntities []model.Channels
	if err := c.gorm.WithContext(ctx).Where(
		`"id" IN ?`,
		uniqueChannelsIds,
	).Find(&channelsEntities).Error; err != nil {
		c.logger.Error("cannot get channels entities", slog.Any("err", err))
		return struct{}{}
	}

	timeoutCtx, _ := context.WithTimeout(context.Background(), 5*time.Minute)

	for _, r := range req {
		var channelEntity *model.Channels
		for _, channel := range channelsEntities {
			if channel.ID == r.ChannelID {
				channelEntity = &channel
				break
			}
		}

		if channelEntity == nil {
			continue
		}

		if r.ChannelID == r.UserID || channelEntity.BotID == r.UserID {
			continue
		}

		go func() {
			if err := c.twitchActions.Ban(
				timeoutCtx,
				twitchactions.BanOpts{
					Duration:       r.BanTime,
					Reason:         r.Reason,
					BroadcasterID:  r.ChannelID,
					UserID:         r.UserID,
					ModeratorID:    channelEntity.BotID,
					IsModerator:    r.IsModerator,
					AddModAfterBan: r.AddModAfterBan,
				},
			); err != nil {
				c.logger.Error("cannot ban user", slog.Any("err", err))
			}
		}()
	}

	return struct{}{}
}
