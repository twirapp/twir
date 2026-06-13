package mappers

import (
	"unicode/utf8"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/entities/platform"
)

func EventSubChatMessageToBus(event eventsub.ChannelChatMessageEvent) generic.ChatMessage {
	fragments := make([]generic.ChatMessageMessageFragment, 0, len(event.Message.Fragments))

	startFragmentPosition := 0
	for _, fragment := range event.Message.Fragments {
		var cheerMote *generic.ChatMessageMessageFragmentCheermote
		var emote *generic.ChatMessageMessageFragmentEmote
		var mention *generic.ChatMessageMessageFragmentMention

		if fragment.Cheermote != nil {
			cheerMote = &generic.ChatMessageMessageFragmentCheermote{
				Prefix: fragment.Cheermote.Prefix,
				Bits:   int64(fragment.Cheermote.Bits),
				Tier:   int64(fragment.Cheermote.Tier),
			}
		}

		if fragment.Emote != nil {
			formats := make([]string, 0, len(fragment.Emote.Format))
			for _, f := range fragment.Emote.Format {
				formats = append(formats, string(f))
			}

			emote = &generic.ChatMessageMessageFragmentEmote{
				ID:         fragment.Emote.Id,
				EmoteSetID: fragment.Emote.EmoteSetId,
				OwnerID:    fragment.Emote.OwnerId,
				Format:     formats,
			}
		}

		if fragment.Mention != nil {
			mention = &generic.ChatMessageMessageFragmentMention{
				UserID:    fragment.Mention.UserId,
				UserName:  fragment.Mention.UserName,
				UserLogin: fragment.Mention.UserLogin,
			}
		}

		position := generic.ChatMessageMessageFragmentPosition{
			Start: startFragmentPosition,
			End:   startFragmentPosition + utf8.RuneCountInString(fragment.Text),
		}

		fragments = append(
			fragments,
			generic.ChatMessageMessageFragment{
				Type:      convertFragmentTypeToEnumValue(string(fragment.Type)),
				Text:      fragment.Text,
				Cheermote: cheerMote,
				Emote:     emote,
				Mention:   mention,
				Position:  position,
			},
		)

		startFragmentPosition += utf8.RuneCountInString(fragment.Text)
	}

	badges := make([]generic.ChatMessageBadge, 0, len(event.Badges))
	for _, badge := range event.Badges {
		badges = append(
			badges,
			generic.ChatMessageBadge{
				ID:    badge.Id,
				SetID: badge.SetId,
				Info:  badge.Info,
				Text:  badge.Info,
			},
		)
	}

	var cheer *generic.ChatMessageCheer
	if event.Cheer != nil {
		cheer = &generic.ChatMessageCheer{Bits: int64(event.Cheer.Bits)}
	}

	var reply *generic.ChatMessageReply
	if event.Reply != nil {
		reply = &generic.ChatMessageReply{
			ParentMessageId:   event.Reply.ParentMessageId,
			ParentMessageBody: event.Reply.ParentMessageBody,
			ParentUserId:      event.Reply.ParentUserId,
			ParentUserName:    event.Reply.ParentUserName,
			ParentUserLogin:   event.Reply.ParentUserLogin,
			ThreadMessageId:   event.Reply.ThreadMessageId,
			ThreadUserId:      event.Reply.ThreadUserId,
			ThreadUserName:    event.Reply.ThreadUserName,
			ThreadUserLogin:   event.Reply.ThreadUserLogin,
		}
	}

	return generic.ChatMessage{
		ID:                          event.MessageId,
		BroadcasterUserId:           event.BroadcasterUserId,
		BroadcasterUserName:         event.BroadcasterUserName,
		BroadcasterUserLogin:        event.BroadcasterUserLogin,
		ChatterUserId:               event.ChatterUserId,
		ChatterUserName:             event.ChatterUserName,
		ChatterUserLogin:            event.ChatterUserLogin,
		MessageID:                   event.MessageId,
		Message: &generic.ChatMessageMessage{
			Text:      event.Message.Text,
			Fragments: fragments,
		},
		Text:                        event.Message.Text,
		Platform:                    string(platform.PlatformTwitch),
		PlatformChannelID:           event.BroadcasterUserId,
		ChannelID:                   event.BroadcasterUserId,
		UserID:                      event.ChatterUserId,
		SenderID:                    event.ChatterUserId,
		SenderLogin:                 event.ChatterUserLogin,
		SenderDisplayName:           event.ChatterUserName,
		Color:                       event.Color,
		Badges:                      badges,
		MessageType:                 string(event.Type),
		Cheer:                       cheer,
		Reply:                       reply,
		ChannelPointsCustomRewardId: event.ChannelPointsCustomRewardId,
	}
}

func convertFragmentTypeToEnumValue(t string) generic.FragmentType {
	switch t {
	case "text":
		return generic.FragmentType_TEXT
	case "cheermote":
		return generic.FragmentType_CHEERMOTE
	case "emote":
		return generic.FragmentType_EMOTE
	case "mention":
		return generic.FragmentType_MENTION
	default:
		return generic.FragmentType_TEXT
	}
}
