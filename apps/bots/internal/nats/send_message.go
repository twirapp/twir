package nats_handlers

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/nats/bots"
)

func (c *NatsHandlers) SendMessage(m *nats.Msg) {
	defer m.Ack()

	data := bots.SendMessage{}
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

	channelName := data.ChannelName

	if channelName == nil {
		usersReq, err := bot.Api.Client.GetUsers(&helix.UsersParams{
			IDs: []string{data.ChannelId},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(usersReq.Data.Users) == 0 {
			fmt.Println("zero length users on send message")
			return
		}
		channelName = &usersReq.Data.Users[0].Login
	}

	if data.IsAnnounce != nil && *data.IsAnnounce == true {
		bot.Api.Client.SendChatAnnouncement(&helix.SendChatAnnouncementParams{
			BroadcasterID: channel.ID,
			ModeratorID:   channel.BotID,
			Message:       data.Message,
		})
	} else {
		bot.SayWithRateLimiting(*channelName, data.Message, nil)
	}

	return
}
