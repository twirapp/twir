package twitchactions

import (
	"context"
	"errors"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/twitch"
)

type DeleteMessageOpts struct {
	BroadcasterID string
	ModeratorID   string
	MessageID     string
}

func (c *TwitchActions) DeleteMessage(ctx context.Context, opts DeleteMessageOpts) error {
	twitchClient, err := twitch.NewBotClientWithContext(ctx, opts.ModeratorID, c.config, c.twirBus)
	if err != nil {
		return err
	}

	resp, err := twitchClient.DeleteChatMessage(
		&helix.DeleteChatMessageParams{
			BroadcasterID: opts.BroadcasterID,
			ModeratorID:   opts.ModeratorID,
			MessageID:     opts.MessageID,
		},
	)
	if err != nil {
		return err
	}

	if resp.ErrorMessage != "" {
		return errors.New(resp.ErrorMessage)
	}

	return nil
}
