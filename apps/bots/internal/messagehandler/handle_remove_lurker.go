package messagehandler

import (
	"context"
	"log/slog"
	"time"

	model "github.com/twirapp/twir/libs/gomodels"
)

var removeLurkerRedisCacheKey = "cache:bots:remove_lurkers:"

func (c *MessageHandler) handleRemoveLurkerBatched(ctx context.Context, data []handleMessage) {
	for _, msg := range data {
		if exists, err := c.redis.Exists(
			ctx,
			removeLurkerRedisCacheKey+msg.ChatterUserId,
		).Result(); err != nil {
			c.logger.Error("cannot remove lurker", slog.Any("err", err))
			continue
		} else if exists == 1 {
			continue
		}

		ignoredUser := &model.IgnoredUser{}
		err := c.gorm.WithContext(ctx).Where(`"id" = ?`, msg.ChatterUserId).Find(ignoredUser).Error
		if err != nil {
			c.logger.Error("cannot remove lurker", slog.Any("err", err))
			continue
		}

		if ignoredUser.ID != "" && !ignoredUser.Force {
			err = c.gorm.WithContext(ctx).Delete(ignoredUser).Error
			if err != nil {
				c.logger.Error("cannot remove lurker", slog.Any("err", err))
				continue
			}
		}

		err = c.redis.Set(
			ctx,
			removeLurkerRedisCacheKey+msg.ChatterUserId,
			"",
			1*time.Hour,
		).Err()
		if err != nil {
			c.logger.Error("cannot remove lurker", slog.Any("err", err))
			continue
		}
	}
}

func (c *MessageHandler) handleRemoveLurker(_ context.Context, msg handleMessage) error {
	c.messagesLurkersBatcher.Add(msg)
	return nil
}
