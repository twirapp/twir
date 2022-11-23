package nats_handlers

import (
	"time"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/satont/tsuwari/apps/bots/types"
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
		rateLimitedChannel := bot.RateLimiters.Channels.Items[data.UserName]
		if rateLimitedChannel == nil {
			bot.RateLimiters.Channels.Lock()
			defer bot.RateLimiters.Channels.Unlock()
			l, _ := ratelimiting.NewSlidingWindow(2, 30*time.Second)
			bot.RateLimiters.Channels.Items[data.UserName] = &types.Channel{
				Limiter: l,
			}
		}
		bot.Join(data.UserName)
	}
	if data.Action == "part" {
		delete(bot.RateLimiters.Channels.Items, data.UserName)
		bot.Depart(data.UserName)
	}
}
