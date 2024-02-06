package messagehandler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"gorm.io/gorm"
)

func (c *MessageHandler) handleGreetings(ctx context.Context, msg handleMessage) error {
	if msg.DbStream == nil {
		return nil
	}

	entity := model.ChannelsGreetings{}
	err := c.gorm.
		WithContext(ctx).
		Where(
			`"channelId" = ? AND "userId" = ? AND "processed" = ? AND "enabled" = ?`,
			msg.GetBroadcasterUserId(),
			msg.GetChatterUserId(),
			false,
			true,
		).
		First(&entity).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	go func() {
		alert := model.ChannelAlert{}
		if err := c.gorm.
			WithContext(ctx).
			Where(
				"channel_id = ? AND greetings_ids && ?",
				msg.GetBroadcasterUserId(),
				pq.StringArray{entity.ID},
			).Find(&alert).Error; err != nil {
			c.logger.Error("cannot find channel alert", slog.Any("err", err))
			return
		}

		if alert.ID == "" {
			return
		}

		c.websocketsGrpc.TriggerAlert(
			ctx,
			&websockets.TriggerAlertRequest{
				ChannelId: msg.GetBroadcasterUserId(),
				AlertId:   alert.ID,
			},
		)
	}()

	if err = c.gorm.WithContext(ctx).Model(&entity).Where("id = ?", entity.ID).Select("*").Updates(
		map[string]any{
			"processed": true,
		},
	).Error; err != nil {
		return err
	}

	requestStruct := &parser.ParseTextRequestData{
		Channel: &parser.Channel{
			Id:   msg.GetBroadcasterUserId(),
			Name: msg.GetBroadcasterUserLogin(),
		},
		Message: &parser.Message{
			Text: entity.Text,
			Id:   msg.GetMessageId(),
		},
		Sender: &parser.Sender{
			Id:          msg.GetChatterUserId(),
			Name:        msg.GetChatterUserLogin(),
			DisplayName: msg.GetChatterUserName(),
			Badges:      createUserBadges(msg.GetBadges()),
		},
		ParseVariables: lo.ToPtr(true),
	}

	res, err := c.parserGrpc.ParseTextResponse(context.Background(), requestStruct)
	if err != nil {
		return err
	}

	for _, r := range res.Responses {
		c.twitchActions.SendMessage(
			ctx, twitchactions.SendMessageOpts{
				BroadcasterID:        msg.GetBroadcasterUserId(),
				SenderID:             msg.DbChannel.BotID,
				Message:              r,
				ReplyParentMessageID: lo.If(entity.IsReply, msg.GetMessageId()).Else(""),
			},
		)
	}

	_, err = c.eventsGrpc.GreetingSended(
		context.Background(),
		&events.GreetingSendedMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: msg.GetBroadcasterUserId()},
			UserId:          msg.GetChatterUserId(),
			UserName:        msg.GetChatterUserLogin(),
			UserDisplayName: msg.GetChatterUserName(),
			GreetingText:    entity.Text,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
