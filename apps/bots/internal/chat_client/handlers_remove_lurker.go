package chat_client

import (
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *ChatClient) removeUserFromLurkers(userId string) {
	ignoredUser := &model.IgnoredUser{}
	err := c.services.DB.Where(`"id" = ?`, userId).Find(ignoredUser).Error
	if err != nil {
		c.services.Logger.Error(
			"cannot find lurker",
			slog.Any("err", err),
			slog.String("channelId", userId),
		)
		return
	}

	if ignoredUser.ID != "" {
		err = c.services.DB.Delete(ignoredUser).Error
		if err != nil {
			c.services.Logger.Error(
				"cannot remove lurker",
				slog.Any("err", err),
				slog.String("channelId", userId),
			)
		}
	}
}
