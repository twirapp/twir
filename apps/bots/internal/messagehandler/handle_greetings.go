package messagehandler

import (
	"context"
	"log/slog"

	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/repositories/greetings"
	greetingsmodel "github.com/twirapp/twir/libs/repositories/greetings/model"
)

func (c *MessageHandler) handleGreetings(ctx context.Context, msg handleMessage) error {
	if msg.DbStream == nil {
		return nil
	}

	allGreetings, err := c.greetingsCache.Get(ctx, msg.BroadcasterUserId)
	if err != nil {
		return err
	}

	var greeting *greetingsmodel.Greeting
	for _, g := range allGreetings {
		if g.UserID == msg.ChatterUserId && g.Enabled && !g.Processed {
			greeting = &g
			break
		}
	}

	if greeting == nil {
		return nil
	}

	go func() {
		alert := model.ChannelAlert{}
		if err := c.gorm.
			WithContext(ctx).
			Where(
				"channel_id = ? AND greetings_ids && ?",
				msg.BroadcasterUserId,
				pq.StringArray{greeting.ID.String()},
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

	if _, err := c.greetingsRepository.Update(
		ctx,
		greeting.ID,
		greetings.UpdateInput{
			Processed: lo.ToPtr(true),
		},
	); err != nil {
		return err
	}

	if err = c.greetingsCache.Invalidate(ctx, msg.BroadcasterUserId); err != nil {
		return err
	}

	res, err := c.bus.Parser.ParseVariablesInText.Request(
		ctx, parser.ParseVariablesInTextRequest{
			ChannelID:   msg.BroadcasterUserId,
			ChannelName: msg.BroadcasterUserLogin,
			Text:        greeting.Text,
			UserID:      msg.ChatterUserId,
			UserLogin:   msg.ChatterUserLogin,
			UserName:    msg.ChatterUserName,
		},
	)
	if err != nil {
		return err
	}

	if res.Data.Text != "" {
		c.twitchActions.SendMessage(
			ctx, twitchactions.SendMessageOpts{
				BroadcasterID:        msg.BroadcasterUserId,
				SenderID:             msg.DbChannel.BotID,
				Message:              res.Data.Text,
				ReplyParentMessageID: lo.If(greeting.IsReply, msg.MessageId).Else(""),
			},
		)
	}

	if greeting.WithShoutOut {
		c.twitchActions.ShoutOut(
			ctx,
			twitchactions.ShoutOutInput{
				BroadcasterID: msg.BroadcasterUserId,
				TargetID:      greeting.UserID,
			},
		)
	}

	_, err = c.eventsGrpc.GreetingSended(
		ctx,
		&events.GreetingSendedMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: msg.BroadcasterUserId},
			UserId:          msg.ChatterUserId,
			UserName:        msg.ChatterUserLogin,
			UserDisplayName: msg.ChatterUserName,
			GreetingText:    greeting.Text,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
