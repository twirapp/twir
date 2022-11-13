package nats_handlers

import (
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/satont/tsuwari/libs/nats/bots"
)

func (c *NatsHandlers) JoinOrLeave(m *nats.Msg) {
	defer m.Ack()

	data := bots.JoinOrLeaveRequest{}
	err := proto.Unmarshal(m.Data, &data)
	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	bot, ok := c.botsService.Instances[data.BotId]
	if !ok {
		return
	}

	if data.Action == "join" {
		bot.Join(data.UserName)
	}
	if data.Action == "part" {
		bot.Depart(data.UserName)
	}

	return
}
