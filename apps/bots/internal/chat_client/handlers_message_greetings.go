package chat_client

import (
	"context"
	"log/slog"

	"github.com/lib/pq"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"go.uber.org/zap"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/parser"
)

func (c *ChatClient) handleGreetings(
	msg Message,
	userBadges []string,
) {
	if msg.DbStream.ID == "" {
		return
	}

	entity := model.ChannelsGreetings{}
	err := c.services.DB.
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
		c.services.Logger.Error(
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

	defer func() {
		alert := model.ChannelAlert{}
		if err := c.services.DB.Where(
			"channel_id = ? AND greetings_ids && ?",
			msg.Channel.ID,
			pq.StringArray{entity.ID},
		).Find(&alert).Error; err != nil {
			zap.S().Error(err)
			return
		}

		if alert.ID == "" {
			return
		}

		c.services.WebsocketsGrpc.TriggerAlert(
			context.Background(),
			&websockets.TriggerAlertRequest{
				ChannelId: msg.Channel.ID,
				AlertId:   alert.ID,
			},
		)
	}()

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

	res, err := c.services.ParserGrpc.ParseTextResponse(context.Background(), requestStruct)
	if err != nil {
		c.services.Logger.Error(
			"cannot parse text response of greeting",
			slog.Any("req", requestStruct),
		)
		return
	}

	for _, r := range res.Responses {
		c.Say(
			SayOpts{
				Channel:   msg.Channel.Name,
				Text:      r,
				ReplyTo:   lo.If(entity.IsReply, &msg.ID).Else(nil),
				WithLimit: true,
			},
		)
	}

	if err = c.services.DB.Model(&entity).Where("id = ?", entity.ID).Select("*").Updates(
		map[string]any{
			"processed": true,
		},
	).Error; err != nil {
		c.services.Logger.Error("cannot update greeting", slog.Any("err", err))
	}

	_, err = c.services.EventsGrpc.GreetingSended(
		context.Background(),
		&events.GreetingSendedMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: msg.Channel.ID},
			UserId:          msg.User.ID,
			UserName:        msg.User.Name,
			UserDisplayName: msg.User.DisplayName,
			GreetingText:    entity.Text,
		},
	)
	if err != nil {
		c.services.Logger.Error("cannot send greetings event", slog.Any("err", err))
	}
}
