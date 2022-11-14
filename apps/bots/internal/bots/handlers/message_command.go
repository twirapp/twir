package handlers

import (
	"fmt"
	"time"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/libs/nats/parser"
)

func (c *Handlers) handleCommand(nats *nats.Conn, msg irc.PrivateMessage, userBadges []string) {
	requestStruct := parser.Request{
		Sender: &parser.Sender{
			Id:          msg.User.ID,
			Name:        msg.User.Name,
			DisplayName: msg.User.DisplayName,
			Badges:      userBadges,
		},
		Channel: &parser.Channel{
			Id:   msg.RoomID,
			Name: msg.Channel,
		},
		Message: &parser.Message{
			Id:   msg.ID,
			Text: msg.Message,
		},
	}
	bytes, err := proto.Marshal(&requestStruct)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := nats.Request("parser.handleProcessCommand", bytes, 5*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	responseStruct := parser.Response{}
	if err = proto.Unmarshal(res.Data, &responseStruct); err != nil {
		return
	}

	if responseStruct.KeepOrder != nil && *responseStruct.KeepOrder {
		for _, v := range responseStruct.Responses {
			r := v
			go c.BotClient.SayWithRateLimiting(
				msg.Channel,
				r,
				lo.If(responseStruct.IsReply, lo.ToPtr(msg.ID)).Else(nil),
			)
		}
	} else {
		for _, r := range responseStruct.Responses {
			c.BotClient.SayWithRateLimiting(
				msg.Channel,
				r,
				lo.If(responseStruct.IsReply, lo.ToPtr(msg.ID)).Else(nil),
			)
		}
	}

	return
}
