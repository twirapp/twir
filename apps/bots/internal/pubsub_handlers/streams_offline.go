package pubsub_handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/satont/twir/libs/pubsub"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *handlers) streamsOffline(data []byte) {
	streamOfflineStruct := &pubsub.StreamOfflineMessage{}
	if err := json.Unmarshal(data, &streamOfflineStruct); err != nil {
		c.logger.Error("cannot unmarshal incoming data", slog.Any("err", err))
		return
	}

	channel := model.Channels{}
	if err := c.db.Where("id = ?", streamOfflineStruct.ChannelID).Find(&channel).Error; err != nil {
		c.logger.Error("cannot find channel", slog.String("channelId", streamOfflineStruct.ChannelID))
		return
	}

	if channel.ID == "" {
		return
	}

	//db.Model(&model.ChannelsGreetings{}).
	//	Where(`"channelId" = ?`, channel.ID).
	//	Update("processed", false)
}
