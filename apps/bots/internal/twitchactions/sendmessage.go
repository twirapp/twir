package twitchactions

import (
	"context"
	"fmt"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/libs/twitch"
)

type SendMessageOpts struct {
	BroadcasterID        string
	SenderID             string
	Message              string
	ReplyParentMessageID string
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

	resp, err := twitchClient.SendChatMessage(
		&helix.SendChatMessageParams{
			BroadcasterID:        opts.BroadcasterID,
			SenderID:             opts.SenderID,
			Message:              validateResponseSlashes(text),
			ReplyParentMessageID: opts.ReplyParentMessageID,
		},
	)
	if err != nil {
		return err
	}

	if resp.ErrorMessage != "" {
		return fmt.Errorf("cannot send message: %w", resp.ErrorMessage)
	}

	return nil
}
