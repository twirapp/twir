package messagehandler

import (
	"context"

	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/greetings"
	greetingsmodel "github.com/twirapp/twir/libs/repositories/greetings/model"
	"github.com/twirapp/twir/libs/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *MessageHandler) handleGreetings(ctx context.Context, msg twitch.TwitchChatMessage) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.String("function.name", utils.GetFuncName()))

	if msg.EnrichedData.ChannelStream == nil {
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

	mentions := make(
		[]twitch.ChatMessageMessageFragmentMention,
		0,
		len(msg.Message.Fragments),
	)
	if msg.Message != nil {
		for _, f := range msg.Message.Fragments {
			if f.Type != twitch.FragmentType_MENTION {
				continue
			}
			if f.Mention != nil {
				mentions = append(mentions, *f.Mention)
			}
		}
	}

	res, err := c.twirBus.Parser.ParseVariablesInText.Request(
		ctx, parser.ParseVariablesInTextRequest{
			ChannelID:   msg.BroadcasterUserId,
			ChannelName: msg.BroadcasterUserLogin,
			Text:        greeting.Text,
			UserID:      msg.ChatterUserId,
			UserLogin:   msg.ChatterUserLogin,
			UserName:    msg.ChatterUserName,
			Mentions:    mentions,
		},
	)
	if err != nil {
		return err
	}

	if res.Data.Text != "" {
		c.twitchActions.SendMessage(
			ctx, twitchactions.SendMessageOpts{
				BroadcasterID:        msg.BroadcasterUserId,
				SenderID:             msg.EnrichedData.DbChannel.BotID,
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

	err = c.twirBus.Events.GreetingSended.Publish(
		ctx,
		events.GreetingSendedMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   msg.BroadcasterUserId,
				ChannelName: msg.BroadcasterUserLogin,
			},
			UserID:          msg.ChatterUserId,
			UserName:        msg.ChatterUserLogin,
			UserDisplayName: msg.ChatterUserName,
			GreetingText:    greeting.Text,
		},
	)
	if err != nil {
		return err
	}

	alert := model.ChannelAlert{}
	if err := c.gorm.
		WithContext(ctx).
		Where(
			"channel_id = ? AND greetings_ids && ?",
			msg.BroadcasterUserId,
			pq.StringArray{greeting.ID.String()},
		).Find(&alert).Error; err != nil {
		c.logger.Error("cannot find channel alert", logger.Error(err))
		return err
	}

	if alert.ID != "" {
		c.websocketsGrpc.TriggerAlert(
			ctx,
			&websockets.TriggerAlertRequest{
				ChannelId: msg.BroadcasterUserId,
				AlertId:   alert.ID,
			},
		)
	}

	return nil
}
