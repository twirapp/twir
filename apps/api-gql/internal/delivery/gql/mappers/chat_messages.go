package mappers

import (
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func ChatMessageToGQL(input entity.ChatMessage) gqlmodel.ChatMessage {
	badges := make([]gqlmodel.ChatMessageBadge, 0, len(input.Badges))
	for _, b := range input.Badges {
		badge := gqlmodel.ChatMessageBadge{
			SetID: b.SetID,
		}
		if b.VersionID != "" {
			badge.VersionID = lo.ToPtr(b.VersionID)
		}
		if b.Text != "" {
			badge.Text = lo.ToPtr(b.Text)
		}

		badges = append(badges, badge)
	}

	fragments := make([]gqlmodel.ChatMessageFragment, 0, len(input.Fragments))
	for _, f := range input.Fragments {
		fragment := gqlmodel.ChatMessageFragment{
			Type: f.Type,
			Text: f.Text,
		}
		if f.EmoteID != "" {
			fragment.EmoteID = lo.ToPtr(f.EmoteID)
		}
		if f.EmoteURL != "" {
			fragment.EmoteURL = lo.ToPtr(f.EmoteURL)
		}

		fragments = append(fragments, fragment)
	}

	var messageID *string
	if input.MessageID != "" {
		messageID = lo.ToPtr(input.MessageID)
	}

	var messageType *string
	if input.MessageType != "" {
		messageType = lo.ToPtr(input.MessageType)
	}

	var announceColor *string
	if input.AnnounceColor != "" {
		announceColor = lo.ToPtr(input.AnnounceColor)
	}

	return gqlmodel.ChatMessage{
		ID:              input.ID,
		Platform:        input.Platform,
		ChannelID:       input.ChannelID,
		ChannelName:     input.ChannelName,
		ChannelLogin:    input.ChannelLogin,
		UserID:          input.UserID,
		UserName:        input.UserName,
		UserDisplayName: input.UserDisplayName,
		UserColor:       input.UserColor,
		Text:            input.Text,
		CreatedAt:       input.CreatedAt,
		MessageID:       messageID,
		MessageType:     messageType,
		AnnounceColor:   announceColor,
		Badges:          badges,
		Fragments:       fragments,
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
		event.UserLogin = lo.ToPtr(e.UserLogin)
	}
	if e.MessageID != "" {
		event.DeletedMessageID = lo.ToPtr(e.MessageID)
	}

	return event
}
