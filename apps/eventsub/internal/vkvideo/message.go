package vkvideo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
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
		ChatMessage *apiChatMessage `json:"chat_message"`
	} `json:"data"`
}

type apiChatMessage struct {
	ID        int64                `json:"id"`
	CreatedAt int64                `json:"created_at"`
	Author    apiChatAuthor        `json:"author"`
	Parts     []apiChatMessagePart `json:"parts"`
}

type apiChatAuthor struct {
	ID          int64  `json:"id"`
	Nick        string `json:"nick"`
	IsOwner     bool   `json:"is_owner"`
	IsModerator bool   `json:"is_moderator"`
}

type apiChatMessagePart struct {
	Text *apiChatTextPart `json:"text"`
}

type apiChatTextPart struct {
	Content string `json:"content"`
}

func parseChatPublication(payload []byte) (chatPublication, bool, error) {
	var envelope publicationEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return chatPublication{}, false, fmt.Errorf("decode Centrifugo publication: %w", err)
	}
	if envelope.Type != "channel_chat_message_send" {
		return chatPublication{}, false, nil
	}
	if envelope.Data.ChatMessage == nil {
		return chatPublication{}, false, errors.New("VK Video chat publication does not contain a message")
	}

	message := envelope.Data.ChatMessage
	var text strings.Builder
	for _, part := range message.Parts {
		if part.Text == nil {
			continue
		}
		text.WriteString(part.Text.Content)
	}

	return chatPublication{
		ID:        strconv.FormatInt(message.ID, 10),
		CreatedAt: strconv.FormatInt(message.CreatedAt, 10),
		Author: chatAuthor{
			ID:                 strconv.FormatInt(message.Author.ID, 10),
			Nick:               message.Author.Nick,
			Name:               message.Author.Nick,
			DisplayName:        message.Author.Nick,
			IsOwner:            message.Author.IsOwner,
			IsChatModerator:    message.Author.IsModerator,
			IsChannelModerator: false,
		},
		Text: text.String(),
	}, true, nil
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
		ID:                uuid.NewString(),
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
