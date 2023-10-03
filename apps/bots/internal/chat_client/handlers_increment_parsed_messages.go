package chat_client

import (
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *ChatClient) incrementStreamParsedMessages(channelId string) {
	stream := model.ChannelsStreams{}
	if err := c.services.DB.Where(`"userId" = ?`, channelId).Select(
		"ID",
		"ParsedMessages",
	).Find(&stream).Error; err != nil {
		c.services.Logger.Error(
			"cannot get channel stream",
			slog.Any("err", err),
			slog.String("channelId", channelId),
		)
		return
	}

	if stream.ID != "" {
		if err := c.services.DB.Model(&stream).Update(
			"parsedMessages",
			stream.ParsedMessages+1,
		).Error; err != nil {
			c.services.Logger.Error(
				"cannot increment parsed messages",
				slog.Any("err", err),
				slog.String("channelId", channelId),
			)
		}
	}
}
