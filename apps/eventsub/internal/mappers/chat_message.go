package mappers

import (
	"unicode/utf8"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/twitch"
)

func EventSubChatMessageToBus(event eventsub.ChannelChatMessageEvent) twitch.TwitchChatMessage {
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
			formats := make([]string, 0, len(fragment.Emote.Format))
			for _, f := range fragment.Emote.Format {
				formats = append(formats, string(f))
			}

			emote = &twitch.ChatMessageMessageFragmentEmote{
				Id:         fragment.Emote.Id,
				EmoteSetId: fragment.Emote.EmoteSetId,
				OwnerId:    fragment.Emote.OwnerId,
				Format:     formats,
			}
		}

		if fragment.Mention != nil {
			mention = &twitch.ChatMessageMessageFragmentMention{
				UserId:    fragment.Mention.UserId,
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

	badges := make([]twitch.ChatMessageBadge, 0, len(event.Badges))
	for _, badge := range event.Badges {
		badges = append(
			badges,
			twitch.ChatMessageBadge{
				Id:    badge.Id,
				SetId: badge.SetId,
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

	return twitch.TwitchChatMessage{
		ID:                   event.MessageId,
		BroadcasterUserId:    event.BroadcasterUserId,
		BroadcasterUserName:  event.BroadcasterUserName,
		BroadcasterUserLogin: event.BroadcasterUserLogin,
		ChatterUserId:        event.ChatterUserId,
		ChatterUserName:      event.ChatterUserName,
		ChatterUserLogin:     event.ChatterUserLogin,
		MessageId:            event.MessageId,
		Message: &twitch.ChatMessageMessage{
			Text:      event.Message.Text,
			Fragments: fragments,
		},
		Color:                       event.Color,
		Badges:                      badges,
		MessageType:                 string(event.Type),
		Cheer:                       cheer,
		Reply:                       reply,
		ChannelPointsCustomRewardId: event.ChannelPointsCustomRewardId,
	}
}

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
