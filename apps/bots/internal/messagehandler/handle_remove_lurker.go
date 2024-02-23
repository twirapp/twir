package messagehandler

import (
	"context"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *MessageHandler) handleRemoveLurker(ctx context.Context, msg handleMessage) error {
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

	return nil
}
