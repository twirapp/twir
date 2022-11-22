package nats_handlers

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	model "github.com/satont/tsuwari/libs/gomodels"
)

type streamOnlineData struct {
	StreamID  string `json:"streamId"`
	ChannelID string `json:"channelId"`
}

type streamOnline struct {
	Pattern string           `json:"pattern"`
	Data    streamOnlineData `json:"data"`
}

func (c *NatsHandlers) StreamOnline(m *nats.Msg) {
	data := streamOnline{}
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

	err := c.db.Model(&model.ChannelsGreetings{}).
		Where(`"channelId" = ?`, channel.ID).
		Update("processed", false).Error
	if err != nil {
		c.logger.Sugar().Error(err)
	}
}
