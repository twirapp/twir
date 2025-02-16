package messagehandler

import (
	"context"
	"time"

	model "github.com/satont/twir/libs/gomodels"
)

var removeLurkerRedisCacheKey = "cache:bots:remove_lurkers:"

func (c *MessageHandler) handleRemoveLurker(ctx context.Context, msg handleMessage) error {
	if exists, err := c.redis.Exists(
		ctx,
		removeLurkerRedisCacheKey+msg.ChatterUserId,
	).Result(); err != nil {
		return err
	} else if exists == 1 {
		return nil
	}

	ignoredUser := &model.IgnoredUser{}
	err := c.gorm.WithContext(ctx).Where(`"id" = ?`, msg.ChatterUserId).Find(ignoredUser).Error
	if err != nil {
		return err
	}

	if ignoredUser.ID != "" && !ignoredUser.Force {
		err = c.gorm.WithContext(ctx).Delete(ignoredUser).Error
		if err != nil {
			return err
		}
	}

	err = c.redis.Set(
		ctx,
		removeLurkerRedisCacheKey+msg.ChatterUserId,
		"",
		1*time.Hour,
	).Err()
	if err != nil {
		return err
	}

	return nil
}
