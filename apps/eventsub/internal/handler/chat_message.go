package handler

import (
	"context"
	"log/slog"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/twirapp/twir/libs/grpc/shared"
)

func (c *Handler) handleChannelChatMessage(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelChatMessage,
) {
	fragments := make([]*shared.ChatMessageMessageFragment, 0, len(event.Message.Fragments))

	for _, fragment := range event.Message.Fragments {
		var cheerMote *shared.ChatMessageMessageFragmentCheermote
		var emote *shared.ChatMessageMessageFragmentEmote
		var mention *shared.ChatMessageMessageFragmentMention

		if fragment.Cheermote != nil {
			cheerMote = &shared.ChatMessageMessageFragmentCheermote{
				Prefix: fragment.Cheermote.Prefix,
				Bits:   int64(fragment.Cheermote.Bits),
				Tier:   int64(fragment.Cheermote.Tier),
			}
		}

		if fragment.Emote != nil {
			emote = &shared.ChatMessageMessageFragmentEmote{
				Id:         fragment.Emote.ID,
				EmoteSetId: fragment.Emote.EmoteSetID,
				OwnerId:    fragment.Emote.OwnerID,
				Format:     fragment.Emote.Format,
			}
		}

		if fragment.Mention != nil {
			mention = &shared.ChatMessageMessageFragmentMention{
				UserId:    fragment.Mention.UserID,
				UserName:  fragment.Mention.UserName,
				UserLogin: fragment.Mention.UserLogin,
			}
		}

		fragments = append(
			fragments, &shared.ChatMessageMessageFragment{
				Type:      fragment.Type,
				Text:      fragment.Text,
				Cheermote: cheerMote,
				Emote:     emote,
				Mention:   mention,
			},
		)
	}

	badges := make([]*shared.ChatMessageBadge, 0, len(event.Badges))
	for _, badge := range event.Badges {
		badges = append(
			badges,
			&shared.ChatMessageBadge{
				Id:    badge.ID,
				SetId: badge.SetID,
				Info:  badge.Info,
			},
		)
	}

	var cheer *shared.ChatMessageCheer
	if event.Cheer != nil {
		cheer = &shared.ChatMessageCheer{Bits: int64(event.Cheer.Bits)}
	}

	var reply *shared.ChatMessageReply
	if event.Reply != nil {
		reply = &shared.ChatMessageReply{
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

	_, err := c.botsGrpc.HandleChatMessage(
		context.TODO(),
		&shared.TwitchChatMessage{
			BroadcasterUserId:    event.BroadcasterUserID,
			BroadcasterUserName:  event.BroadcasterUserName,
			BroadcasterUserLogin: event.BroadcasterUserLogin,
			ChatterUserId:        event.ChatterUserID,
			ChatterUserName:      event.BroadcasterUserName,
			ChatterUserLogin:     event.ChatterUserLogin,
			MessageId:            event.MessageID,
			Message: &shared.ChatMessageMessage{
				Text:      event.Message.Text,
				Fragments: fragments,
			},
			Color:                       event.Color,
			Badges:                      badges,
			MessageType:                 event.MessageType,
			Cheer:                       cheer,
			Reply:                       reply,
			ChannelPointsCustomRewardId: event.ChannelPointsCustomRewardID,
		},
	)

	if err != nil {
		c.logger.Error("cannot handle message", slog.Any("err", err))
	}
}
