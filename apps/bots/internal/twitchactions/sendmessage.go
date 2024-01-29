package twitchactions

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/libs/twitch"
)

type SendMessageOpts struct {
	BroadcasterID        string
	SenderID             string
	Message              string
	ReplyParentMessageID string
	IsAnnounce           bool
}

func validateResponseSlashes(response string) string {
	if strings.HasPrefix(response, "/me") || strings.HasPrefix(response, "/announce") {
		return response
	} else if strings.HasPrefix(response, "/") {
		return "Slash commands except /me and /announce is disallowed. This response wont be ever sended."
	} else if strings.HasPrefix(response, ".") {
		return `Message cannot start with "." symbol.`
	} else {
		return response
	}
}

func (c *TwitchActions) SendMessage(ctx context.Context, opts SendMessageOpts) error {
	twitchClient, err := twitch.NewBotClientWithContext(ctx, opts.SenderID, c.Config, c.TokensGrpc)
	if err != nil {
		return err
	}

	text := strings.ReplaceAll(opts.Message, "\n", " ")
	textParts := splitTextByLength(text)

	for i, part := range textParts {
		if i > 2 {
			return nil
		}

		var msgErr error
		var errorMessage string

		if !opts.IsAnnounce {
			resp, err := twitchClient.SendChatMessage(
				&helix.SendChatMessageParams{
					BroadcasterID:        opts.BroadcasterID,
					SenderID:             opts.SenderID,
					Message:              validateResponseSlashes(part),
					ReplyParentMessageID: opts.ReplyParentMessageID,
				},
			)
			msgErr = err
			errorMessage = resp.ErrorMessage
		} else {
			resp, err := twitchClient.SendChatAnnouncement(
				&helix.SendChatAnnouncementParams{
					BroadcasterID: opts.BroadcasterID,
					ModeratorID:   opts.SenderID,
					Message:       validateResponseSlashes(part),
				},
			)
			msgErr = err
			errorMessage = resp.ErrorMessage
		}

		if msgErr != nil {
			return err
		}

		if errorMessage != "" {
			return fmt.Errorf("cannot send message: %w", errorMessage)
		}
	}

	return nil
}

func splitTextByLength(text string) []string {
	var parts []string

	i := 500
	for utf8.RuneCountInString(text) > 0 {
		if utf8.RuneCountInString(text) < 500 {
			parts = append(parts, text)
			break
		}
		runned := []rune(text)
		parts = append(parts, string(runned[:i]))
		text = string(runned[i:])
	}

	return parts
}
