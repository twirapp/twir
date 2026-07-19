package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func ChatOverlayMessageEntityToGQL(m entity.ChatMessage) gqlmodel.ChatOverlayMessage {
	badges := make([]gqlmodel.ChatOverlayMessageBadge, 0, len(m.Badges))
	for _, b := range m.Badges {
		badge := gqlmodel.ChatOverlayMessageBadge{
			SetID: b.SetID,
		}
		if b.VersionID != "" {
			badge.VersionID = &b.VersionID
		}
		if b.Text != "" {
			badge.Text = &b.Text
		}

		badges = append(badges, badge)
	}

	fragments := make([]gqlmodel.ChatOverlayMessageFragment, 0, len(m.Fragments))
	for _, f := range m.Fragments {
		fragment := gqlmodel.ChatOverlayMessageFragment{
			Type: f.Type,
			Text: f.Text,
		}
		if f.EmoteID != "" {
			fragment.EmoteID = &f.EmoteID
		}
		if f.EmoteURL != "" {
			fragment.EmoteURL = &f.EmoteURL
		}

		fragments = append(fragments, fragment)
	}

	var senderColor *string
	if m.UserColor != "" {
		senderColor = &m.UserColor
	}

	var announceColor *string
	if m.AnnounceColor != "" {
		announceColor = &m.AnnounceColor
	}

	return gqlmodel.ChatOverlayMessage{
		ID:                m.ID.String(),
		Platform:          m.Platform,
		MessageID:         m.MessageID,
		MessageType:       m.MessageType,
		SenderID:          m.SenderID,
		SenderLogin:       m.UserName,
		SenderDisplayName: m.UserDisplayName,
		SenderColor:       senderColor,
		AnnounceColor:     announceColor,
		Badges:            badges,
		Fragments:         fragments,
		CreatedAt:         m.CreatedAt,
	}
}

func ChatOverlayModerationEventEntityToGQL(
	e entity.ChatOverlayModerationEvent,
) gqlmodel.ChatOverlayModerationEvent {
	event := gqlmodel.ChatOverlayModerationEvent{
		Type:     gqlmodel.ChatOverlayModerationEventType(e.Type),
		Platform: e.Platform,
	}
	if e.UserLogin != "" {
		event.UserLogin = &e.UserLogin
	}
	if e.MessageID != "" {
		event.DeletedMessageID = &e.MessageID
	}

	return event
}
