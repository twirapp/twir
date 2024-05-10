package handler

import (
	"context"
	"log/slog"
	"unicode/utf8"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/grpc/events"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
	"go.opentelemetry.io/otel/attribute"
)

func convertFragmentTypeToEnumValue(t string) twitch.FragmentType {
	switch t {
	case "text":
		return twitch.FragmentType_TEXT
	case "cheermote":
		return twitch.FragmentType_CHEERMOTE
	case "emote":
		return twitch.FragmentType_EMOTE
	case "mention":
		return twitch.FragmentType_MENTION
	default:
		return twitch.FragmentType_TEXT
	}
}

func (c *Handler) handleChannelChatMessage(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelChatMessage,
) {
	ctx := context.Background()
	_, span := c.tracer.Start(ctx, "handleChannelChatMessage")
	span.SetAttributes(
		attribute.String("message_id", event.MessageID),
		attribute.String("channel_id", event.BroadcasterUserID),
	)
	defer span.End()

	fragments := make([]twitch.ChatMessageMessageFragment, 0, len(event.Message.Fragments))

	startFragmentPosition := 0
	for _, fragment := range event.Message.Fragments {
		var cheerMote *twitch.ChatMessageMessageFragmentCheermote
		var emote *twitch.ChatMessageMessageFragmentEmote
		var mention *twitch.ChatMessageMessageFragmentMention

		if fragment.Cheermote != nil {
			cheerMote = &twitch.ChatMessageMessageFragmentCheermote{
				Prefix: fragment.Cheermote.Prefix,
				Bits:   int64(fragment.Cheermote.Bits),
				Tier:   int64(fragment.Cheermote.Tier),
			}
		}

		if fragment.Emote != nil {
			emote = &twitch.ChatMessageMessageFragmentEmote{
				Id:         fragment.Emote.ID,
				EmoteSetId: fragment.Emote.EmoteSetID,
				OwnerId:    fragment.Emote.OwnerID,
				Format:     fragment.Emote.Format,
			}
		}

		if fragment.Mention != nil {
			mention = &twitch.ChatMessageMessageFragmentMention{
				UserId:    fragment.Mention.UserID,
				UserName:  fragment.Mention.UserName,
				UserLogin: fragment.Mention.UserLogin,
			}
		}

		position := twitch.ChatMessageMessageFragmentPosition{
			Start: startFragmentPosition,
			End:   startFragmentPosition + utf8.RuneCountInString(fragment.Text),
		}

		fragments = append(
			fragments,
			twitch.ChatMessageMessageFragment{
				Type:      convertFragmentTypeToEnumValue(fragment.Type),
				Text:      fragment.Text,
				Cheermote: cheerMote,
				Emote:     emote,
				Mention:   mention,
				Position:  position,
			},
		)

		startFragmentPosition += utf8.RuneCountInString(fragment.Text)
	}

	badges := make([]twitch.ChatMessageBadge, 0, len(event.Badges))
	for _, badge := range event.Badges {
		badges = append(
			badges,
			twitch.ChatMessageBadge{
				Id:    badge.ID,
				SetId: badge.SetID,
				Info:  badge.Info,
			},
		)
	}

	var cheer *twitch.ChatMessageCheer
	if event.Cheer != nil {
		cheer = &twitch.ChatMessageCheer{Bits: int64(event.Cheer.Bits)}
	}

	var reply *twitch.ChatMessageReply
	if event.Reply != nil {
		reply = &twitch.ChatMessageReply{
			ParentMessageId:   event.Reply.ParentMessageID,
			ParentMessageBody: event.Reply.ParentMessageBody,
			ParentUserId:      event.Reply.ParentUserID,
			ParentUserName:    event.Reply.ParentUserName,
			ParentUserLogin:   event.Reply.ParentUserLogin,
			ThreadMessageId:   event.Reply.ThreadMessageID,
			ThreadUserId:      event.Reply.ThreadUserID,
			ThreadUserName:    event.Reply.ThreadUserName,
			ThreadUserLogin:   event.Reply.ThreadUserLogin,
		}
	}

	data := twitch.TwitchChatMessage{
		BroadcasterUserId:    event.BroadcasterUserID,
		BroadcasterUserName:  event.BroadcasterUserName,
		BroadcasterUserLogin: event.BroadcasterUserLogin,
		ChatterUserId:        event.ChatterUserID,
		ChatterUserName:      event.ChatterUserName,
		ChatterUserLogin:     event.ChatterUserLogin,
		MessageId:            event.MessageID,
		Message: &twitch.ChatMessageMessage{
			Text:      event.Message.Text,
			Fragments: fragments,
		},
		Color:                       event.Color,
		Badges:                      badges,
		MessageType:                 event.MessageType,
		Cheer:                       cheer,
		Reply:                       reply,
		ChannelPointsCustomRewardId: event.ChannelPointsCustomRewardID,
	}

	if err := c.bus.Bots.ProcessMessage.Publish(data); err != nil {
		c.logger.Error("cannot handle message", slog.Any("err", err))
	}

	channel := &model.Channels{}
	if err := c.gorm.WithContext(ctx).Where(
		"id = ?",
		data.BroadcasterUserId,
	).Find(channel).Error; err != nil {
		c.logger.Error("cannot get channel", slog.Any("err", err))
		return
	}
	if channel.ID == "" {
		return
	}
	if channel.BotID == data.ChatterUserId && c.config.AppEnv == "production" {
		return
	}

	if data.Message.Text[0] == '!' {
		if err := c.bus.Parser.ProcessMessageAsCommand.Publish(data); err != nil {
			c.logger.Error("cannot process command", slog.Any("err", err))
		}
	}
}

func (c *Handler) handleChannelChatMessageDelete(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelChatMessageDelete,
) {
	c.logger.Info(
		"message delete",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserName),
		slog.String("userId", event.TargetUserID),
		slog.String("userName", event.TargetUserLogin),
	)

	if _, err := c.eventsGrpc.ChannelMessageDelete(
		context.Background(),
		&events.ChannelMessageDeleteMessage{
			BaseInfo:             &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			UserId:               event.TargetUserID,
			UserName:             event.TargetUserLogin,
			UserLogin:            event.TargetUserName,
			BroadcasterUserName:  event.BroadcasterUserName,
			BroadcasterUserLogin: event.BroadcasterUserLogin,
			MessageId:            event.MessageID,
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
