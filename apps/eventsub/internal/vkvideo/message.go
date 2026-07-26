package vkvideo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	user_creator "github.com/twirapp/twir/apps/eventsub/internal/services/user-creator"
	"github.com/twirapp/twir/libs/bus-core/generic"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type chatPublication struct {
	ID        string
	CreatedAt string
	Author    chatAuthor
	Text      string
}

type chatAuthor struct {
	ID                 string `json:"id"`
	Nick               string `json:"nick"`
	Name               string `json:"name"`
	DisplayName        string `json:"displayName"`
	AvatarURL          string `json:"avatarUrl"`
	IsOwner            bool   `json:"isOwner"`
	IsChatModerator    bool   `json:"isChatModerator"`
	IsChannelModerator bool   `json:"isChannelModerator"`
}

type publicationEnvelope struct {
	Type string `json:"type"`
	Data struct {
		ID        string         `json:"id"`
		CreatedAt string         `json:"createdAt"`
		Author    chatAuthor     `json:"author"`
		Data      []chatFragment `json:"data"`
	} `json:"data"`
}

type chatFragment struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func parseChatPublication(payload []byte) (chatPublication, bool, error) {
	var envelope publicationEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return chatPublication{}, false, fmt.Errorf("decode Centrifugo publication: %w", err)
	}
	if envelope.Type != "message" {
		return chatPublication{}, false, nil
	}
	if envelope.Data.ID == "" {
		return chatPublication{}, false, errors.New("VK Video chat message id is empty")
	}

	var text strings.Builder
	for _, fragment := range envelope.Data.Data {
		if fragment.Type == "BLOCK_END" || fragment.Content == "" {
			continue
		}
		fragmentText, err := decodeTextFragment(fragment.Content)
		if err != nil {
			return chatPublication{}, false, err
		}
		text.WriteString(fragmentText)
	}

	return chatPublication{
		ID:        envelope.Data.ID,
		CreatedAt: envelope.Data.CreatedAt,
		Author:    envelope.Data.Author,
		Text:      text.String(),
	}, true, nil
}

func decodeTextFragment(content string) (string, error) {
	var parts []json.RawMessage
	if err := json.Unmarshal([]byte(content), &parts); err != nil {
		return "", fmt.Errorf("decode VK Video chat text fragment: %w", err)
	}
	if len(parts) != 3 {
		return "", errors.New("VK Video chat text fragment has unexpected shape")
	}

	var text, style string
	var attributes []json.RawMessage
	if err := json.Unmarshal(parts[0], &text); err != nil {
		return "", fmt.Errorf("decode VK Video chat text: %w", err)
	}
	if err := json.Unmarshal(parts[1], &style); err != nil {
		return "", fmt.Errorf("decode VK Video chat text style: %w", err)
	}
	if err := json.Unmarshal(parts[2], &attributes); err != nil {
		return "", fmt.Errorf("decode VK Video chat text attributes: %w", err)
	}
	if style != "unstyled" || len(attributes) != 0 {
		return "", errors.New("VK Video chat text fragment has unsupported formatting")
	}

	return text, nil
}

func (t *Transport) handlePublication(
	ctx context.Context,
	binding channelplatformentity.ChannelPlatform,
	payload []byte,
) error {
	publication, isChatMessage, err := parseChatPublication(payload)
	if err != nil || !isChatMessage {
		return err
	}

	claimed, err := t.deduplicator.Claim(ctx, publication.ID)
	if err != nil {
		return err
	}
	if !claimed {
		return nil
	}

	sender, err := t.resolveUser(ctx, binding, publication.Author)
	if err != nil {
		return err
	}
	if binding.BotUserID != nil && *binding.BotUserID == sender.ID && *binding.BotUserID != binding.UserID {
		return nil
	}

	message := normalizeChatMessage(binding, publication, sender)
	if err := t.chatMessages.Publish(ctx, message); err != nil {
		return fmt.Errorf("publish VK Video chat message: %w", err)
	}
	if err := t.commands.Publish(ctx, message); err != nil {
		return fmt.Errorf("publish VK Video parser message: %w", err)
	}

	return nil
}

func (t *Transport) resolveUser(
	ctx context.Context,
	binding channelplatformentity.ChannelPlatform,
	author chatAuthor,
) (usersmodel.User, error) {
	displayName := author.DisplayName
	if displayName == "" {
		displayName = author.Nick
	}
	channelID := binding.ChannelID.String()
	user, _, err := t.userCreator.UnsureUser(ctx, user_creator.CreateUserInput{
		UserID:            author.ID,
		PlatformID:        author.ID,
		Platform:          platformentity.PlatformVKVideoLive,
		Login:             author.Nick,
		DisplayName:       displayName,
		ChannelID:         &channelID,
		ShouldUpdateStats: false,
		IsBroadcaster:     author.IsOwner,
		IsModerator:       author.IsChatModerator || author.IsChannelModerator,
	})
	if err != nil {
		return usersmodel.User{}, fmt.Errorf("ensure VK Video chat user: %w", err)
	}
	return *user, nil
}

func normalizeChatMessage(
	binding channelplatformentity.ChannelPlatform,
	publication chatPublication,
	sender usersmodel.User,
) generic.ChatMessage {
	textLength := utf8.RuneCountInString(publication.Text)
	return generic.ChatMessage{
		Message: &generic.ChatMessageMessage{
			Text: publication.Text,
			Fragments: []generic.ChatMessageMessageFragment{{
				Type: generic.FragmentType_TEXT,
				Text: publication.Text,
				Position: generic.ChatMessageMessageFragmentPosition{
					Start: 0,
					End:   textLength,
				},
			}},
		},
		ID:                publication.ID,
		BroadcasterUserId: binding.PlatformChannelID,
		ChatterUserId:     publication.Author.ID,
		ChatterUserName:   publication.Author.Name,
		ChatterUserLogin:  publication.Author.Nick,
		MessageType:       "text",
		Platform:          string(platformentity.PlatformVKVideoLive),
		ChannelID:         binding.ChannelID.String(),
		ChannelBindingID:  binding.ID.String(),
		UserID:            sender.ID.String(),
		PlatformChannelID: binding.PlatformChannelID,
		SenderID:          publication.Author.ID,
		SenderLogin:       publication.Author.Nick,
		SenderDisplayName: publication.Author.DisplayName,
		MessageID:         publication.ID,
		Text:              publication.Text,
		IsBroadcaster:     publication.Author.IsOwner,
		IsModerator:       publication.Author.IsChatModerator || publication.Author.IsChannelModerator,
	}
}
