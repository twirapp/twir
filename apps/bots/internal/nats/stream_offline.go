package nats_handlers

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	model "github.com/satont/tsuwari/libs/gomodels"
)

type streamOfflineData struct {
	ChannelID string `json:"channelId"`
}

type streamOffline struct {
	Pattern string           `json:"pattern"`
	Data    streamOnlineData `json:"data"`
}

func (c *NatsHandlers) StreamOffline(m *nats.Msg) {
	data := streamOffline{}
	if err := json.Unmarshal(m.Data, &data); err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	channel := model.Channels{}
	if err := c.db.Where("id = ?", data.Data.ChannelID).Find(&channel).Error; err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	if channel.ID == "" {
		return
	}

	c.db.Model(&model.ChannelsGreetings{}).
		Where(`"channelId" = ?`, channel.ID).
		Update("processed", false)
}
