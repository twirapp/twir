package handlers

import (
	"context"
	"log/slog"

	"github.com/satont/twir/libs/grpc/generated/events"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/parser"
)

func (c *Handlers) handleGreetings(
	msg *Message,
	userBadges []string,
) {
	stream := &model.ChannelsStreams{}
	err := c.db.Where(`"userId" = ?`, msg.Channel.ID).Find(&stream).Error
	if err != nil {
		c.logger.Error("cannot get stream", slog.Any("err", err))
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
		Find(&entity).
		Error
	if err != nil {
		c.logger.Error(
			"cannot get greeting",
			slog.Any("err", err),
			slog.String("channelId", msg.Channel.ID),
			slog.String("userId", msg.User.ID),
		)
		return
	}

	if entity.ID == "" {
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
		c.logger.Error("cannot parse text response of greeting", slog.Any("req", requestStruct))
		return
	}

	for _, r := range res.Responses {
		validateResponseErr := ValidateResponseSlashes(r)
		if validateResponseErr != nil {
			c.BotClient.SayWithRateLimiting(
				msg.Channel.Name,
				validateResponseErr.Error(),
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

	if err = c.db.Model(&entity).Where("id = ?", entity.ID).Select("*").Updates(
		map[string]any{
			"processed": true,
		},
	).Error; err != nil {
		c.logger.Error("cannot update greeting", slog.Any("err", err))
	}

	_, err = c.eventsGrpc.GreetingSended(
		context.Background(), &events.GreetingSendedMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: msg.Channel.ID},
			UserId:          msg.User.ID,
			UserName:        msg.User.Name,
			UserDisplayName: msg.User.DisplayName,
			GreetingText:    entity.Text,
		},
	)
	c.logger.Error("cannot send greetings event", slog.Any("err", err))
}
