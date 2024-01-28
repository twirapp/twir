package messagehandler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/lib/pq"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"gorm.io/gorm"
)

func (c *MessageHandler) handleGreetings(ctx context.Context, msg handleMessage) error {
	if msg.DbStream == nil {
		return nil
	}

	entity := model.ChannelsGreetings{}
	err := c.gorm.
		WithContext(ctx).
		Where(
			`"channelId" = ? AND "userId" = ? AND "processed" = ? AND "enabled" = ?`,
			msg.GetBroadcasterUserId(),
			msg.GetChatterUserId(),
			false,
			true,
		).
		First(&entity).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	c.pool.Submit(
		func() {
			alert := model.ChannelAlert{}
			if err := c.gorm.
				WithContext(ctx).
				Where(
					"channel_id = ? AND greetings_ids && ?",
					msg.GetBroadcasterUserId(),
					pq.StringArray{entity.ID},
				).Find(&alert).Error; err != nil {
				c.logger.Error("cannot find channel alert", slog.Any("err", err))
				return
			}

			if alert.ID == "" {
				return
			}

			c.websocketsGrpc.TriggerAlert(
				ctx,
				&websockets.TriggerAlertRequest{
					ChannelId: msg.GetBroadcasterUserId(),
					AlertId:   alert.ID,
				},
			)
		},
	)

	return nil
}
