package handlers

import (
	model "github.com/satont/twir/libs/gomodels"
	"log/slog"
)

func (c *Handlers) removeUserFromLurkers(userId string) {
	ignoredUser := &model.IgnoredUser{}
	err := c.db.Where(`"id" = ?`, userId).Find(ignoredUser).Error
	if err != nil {
		c.logger.Error(
			"cannot find lurker",
			slog.Any("err", err),
			slog.String("channelId", userId),
		)
		return
	}

	if ignoredUser.ID != "" {
		err = c.db.Delete(ignoredUser).Error
		if err != nil {
			c.logger.Error(
				"cannot remove lurker",
				slog.Any("err", err),
				slog.String("channelId", userId),
			)
		}
	}
}
