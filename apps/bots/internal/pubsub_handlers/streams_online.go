package pubsub_handlers

import (
	"encoding/json"
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/pubsub"
)

func (c *handlers) streamsOnline(data []byte) {
	streamOnlineStruct := &pubsub.StreamOnlineMessage{}
	if err := json.Unmarshal(data, &streamOnlineStruct); err != nil {
		c.logger.Error("cannot unmarshal incoming data", slog.Any("err", err))
		return
	}

	channel := model.Channels{}
	if err := c.db.Where("id = ?", streamOnlineStruct.ChannelID).Find(&channel).Error; err != nil {
		c.logger.Error("cannot find channel", slog.String("channelId", streamOnlineStruct.ChannelID))
		return
	}

	if channel.ID == "" {
		return
	}

	err := c.db.Model(&model.ChannelsGreetings{}).
		Where(`"channelId" = ?`, channel.ID).
		Update("processed", false).Error
	if err != nil {
		c.logger.Error(
			"cannot update channel greetings", slog.String("channelId", streamOnlineStruct.ChannelID),
			slog.Any("err", err),
		)
	}
}
