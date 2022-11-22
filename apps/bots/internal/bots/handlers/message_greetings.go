package handlers

import (
	"fmt"
	"time"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/golang/protobuf/proto"
	"github.com/samber/lo"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/nats/parser"
	"gorm.io/gorm"
)

func (c *Handlers) handleGreetings(
	msg irc.PrivateMessage,
	userBadges []string,
) {
	entity := model.ChannelsGreetings{}
	err := c.db.
		Where(
			`"channelId" = ? AND "userId" = ? AND "processed" = ? AND "enabled" = ?`,
			msg.RoomID,
			msg.User.ID,
			false,
			true,
		).
		First(&entity).
		Error

	if err != nil && err == gorm.ErrRecordNotFound {
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	requestStruct := parser.ParseResponseRequest{
		Channel: &parser.Channel{
			Id:   msg.RoomID,
			Name: msg.Channel,
		},
		Message: &parser.Message{
			Text: entity.Text,
			Id:   msg.ID,
		},
		Sender: &parser.Sender{
			Id:          msg.User.ID,
			Name:        msg.User.Name,
			DisplayName: msg.User.DisplayName,
			Badges:      userBadges,
		},
		ParseVariables: lo.ToPtr(true),
	}
	bytes, err := proto.Marshal(&requestStruct)
	if err != nil {
		fmt.Println("Parser not answered on request greetings.")
		fmt.Printf("%+v\n", &requestStruct)
		fmt.Println(err)
		return
	}

	responseStruct := parser.ParseResponseResponse{}
	res, err := c.nats.Request("parser.parseTextResponse", bytes, 5*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := proto.Unmarshal(res.Data, &responseStruct); err != nil {
		fmt.Println(err)
		return
	}

	for _, r := range responseStruct.Responses {
		validateResposeErr := ValidateResponseSlashes(r)
		if validateResposeErr != nil {
			c.BotClient.SayWithRateLimiting(
				msg.Channel,
				validateResposeErr.Error(),
				nil,
			)
		} else {
			c.BotClient.SayWithRateLimiting(
				msg.Channel,
				r,
				lo.If(entity.IsReply, &msg.ID).Else(nil),
			)
		}
	}

	c.db.Model(&entity).Where("id = ?", entity.ID).Select("*").Updates(map[string]any{
		"processed": true,
	})
}
