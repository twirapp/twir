package nats_handlers

import (
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/satont/tsuwari/libs/nats/bots"
)

func (c *NatsHandlers) JoinOrLeave(m *nats.Msg) {
	defer m.Ack()

	data := bots.DeleteMessagesRequest{}
	err := proto.Unmarshal(m.Data, &data)
	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	channel := model.Channels{}
	err = c.db.Where("id = ?", data.ChannelId).Find(&channel).Error
	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	if channel.ID == "" {
		return
	}

	bot, ok := c.botsService.Instances[channel.BotID]
	if !ok {
		return
	}

	for _, m := range data.MessageIds {
		bot.Api.Client.DeleteMessage(&helix.DeleteMessageParams{
			BroadcasterID: channel.ID,
			ModeratorID:   channel.BotID,
			MessageID:     m,
		})
	}

	return
}
