package handlers

import (
	"log/slog"
	"strings"
	"time"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *Handlers) storeMessage(
	messageId, channelId, userId, userName, text string,
	canBeDeleted bool,
) {
	entity := model.ChannelChatMessage{
		MessageId:    messageId,
		ChannelId:    channelId,
		UserId:       userId,
		UserName:     userName,
		Text:         strings.ToLower(text),
		CanBeDeleted: canBeDeleted,
		CreatedAt:    time.Now().UTC(),
	}

	err := c.db.Create(&entity).Error
	if err != nil {
		c.logger.Error(
			"cannot save user message to db",
			slog.String("channelId", channelId),
			slog.Group("user", slog.String("id", userId), slog.String("name", userName)),
		)
	}
}
