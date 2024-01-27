package twitchactions

import (
	"context"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/libs/twitch"
)

type SendMessageOpts struct {
	BroadcasterID        string
	SenderID             string
	Message              string
	ReplyParentMessageID string
}

func (c *TwitchActions) SendMessage(ctx context.Context, opts SendMessageOpts) error {
	twitchClient, err := twitch.NewBotClientWithContext(ctx, opts.SenderID, c.Config, c.TokensGrpc)
	if err != nil {
		return err
	}

	resp, err := twitchClient.SendChatMessage(
		&helix.SendChatMessageParams{
			BroadcasterID:        opts.BroadcasterID,
			SenderID:             opts.SenderID,
			Message:              opts.Message,
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
