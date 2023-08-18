package handlers

import (
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *Handlers) incrementStreamParsedMessages(channelId string) {
	stream := model.ChannelsStreams{}
	if err := c.db.Where(`"userId" = ?`, channelId).Select(
		"ID",
		"ParsedMessages",
	).Find(&stream).Error; err != nil {
		c.logger.Error(
			"cannot get channel stream",
			slog.Any("err", err),
			slog.String("channelId", channelId),
		)
		return
	}

	if stream.ID != "" {
		if err := c.db.Model(&stream).Update(
			"parsedMessages",
			stream.ParsedMessages+1,
		).Error; err != nil {
			c.logger.Error(
				"cannot increment parsed messages",
				slog.Any("err", err),
				slog.String("channelId", channelId),
			)
		}
	}
}
