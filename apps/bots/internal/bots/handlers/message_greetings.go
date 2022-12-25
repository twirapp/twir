package handlers

import (
	"context"
	"fmt"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/samber/lo"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
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

	defer c.greetingsCounter.Inc()

	requestStruct := &parser.ParseTextRequestData{
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

	res, err := c.parserGrpc.ParseTextResponse(context.Background(), requestStruct)
	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	for _, r := range res.Responses {
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
