package chat_client

import (
	"log/slog"
	"strings"
	"time"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *ChatClient) storeMessage(
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

	err := c.services.DB.Create(&entity).Error
	if err != nil {
		c.services.Logger.Error(
			"cannot save user message to db",
			slog.String("channelId", channelId),
			slog.Group("user", slog.String("id", userId), slog.String("name", userName)),
		)
	}
}
