package messagehandler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/grpc/events"
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
			msg.BroadcasterUserId,
			msg.ChatterUserId,
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
				msg.BroadcasterUserId,
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
				ChannelId: msg.BroadcasterUserId,
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

	res, err := c.bus.Parser.ParseVariablesInText.Request(
		ctx, parser.ParseVariablesInTextRequest{
			ChannelID:   msg.BroadcasterUserId,
			ChannelName: msg.BroadcasterUserLogin,
			Text:        entity.Text,
			UserID:      msg.ChatterUserId,
			UserLogin:   msg.ChatterUserLogin,
			UserName:    msg.ChatterUserName,
		},
	)
	if err != nil {
		return err
	}

	c.twitchActions.SendMessage(
		ctx, twitchactions.SendMessageOpts{
			BroadcasterID:        msg.BroadcasterUserId,
			SenderID:             msg.DbChannel.BotID,
			Message:              res.Data.Text,
			ReplyParentMessageID: lo.If(entity.IsReply, msg.MessageId).Else(""),
		},
	)

	_, err = c.eventsGrpc.GreetingSended(
		context.Background(),
		&events.GreetingSendedMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: msg.BroadcasterUserId},
			UserId:          msg.ChatterUserId,
			UserName:        msg.ChatterUserLogin,
			UserDisplayName: msg.ChatterUserName,
			GreetingText:    entity.Text,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
