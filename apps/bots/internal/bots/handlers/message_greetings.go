package handlers

import (
	"context"
	"fmt"
	"github.com/satont/tsuwari/libs/grpc/generated/events"

	"github.com/samber/lo"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"gorm.io/gorm"
)

func (c *Handlers) handleGreetings(
	msg *Message,
	userBadges []string,
) {
	stream := &model.ChannelsStreams{}
	err := c.db.Where(`"userId" = ?`, msg.Channel.ID).Find(&stream).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	if stream.ID == "" {
		return
	}

	entity := model.ChannelsGreetings{}
	err = c.db.
		Where(
			`"channelId" = ? AND "userId" = ? AND "processed" = ? AND "enabled" = ?`,
			msg.Channel.ID,
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

	defer c.greetingsCounter.Inc()

	requestStruct := &parser.ParseTextRequestData{
		Channel: &parser.Channel{
			Id:   msg.Channel.ID,
			Name: msg.Channel.Name,
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

	res, err := c.parserGrpc.ParseTextResponse(context.Background(), requestStruct)
	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	for _, r := range res.Responses {
		validateResposeErr := ValidateResponseSlashes(r)
		if validateResposeErr != nil {
			c.BotClient.SayWithRateLimiting(
				msg.Channel.Name,
				validateResposeErr.Error(),
				nil,
			)
		} else {
			c.BotClient.SayWithRateLimiting(
				msg.Channel.Name,
				r,
				lo.If(entity.IsReply, &msg.ID).Else(nil),
			)
		}
	}

	c.db.Model(&entity).Where("id = ?", entity.ID).Select("*").Updates(map[string]any{
		"processed": true,
	})

	c.eventsGrpc.GreetingSended(context.Background(), &events.GreetingSendedMessage{
		BaseInfo:        &events.BaseInfo{ChannelId: msg.Channel.ID},
		UserId:          msg.User.ID,
		UserName:        msg.User.Name,
		UserDisplayName: msg.User.DisplayName,
		GreetingText:    entity.Text,
	})
}
